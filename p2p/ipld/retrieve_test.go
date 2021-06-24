package ipld

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ipfs/go-ipfs/core/coreapi"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	format "github.com/ipfs/go-ipld-format"
	"github.com/lazyledger/lazyledger-core/ipfs"
	"github.com/lazyledger/lazyledger-core/ipfs/plugin"
	"github.com/lazyledger/lazyledger-core/libs/log"
	"github.com/lazyledger/lazyledger-core/types"
	"github.com/lazyledger/lazyledger-core/types/consts"
	"github.com/lazyledger/nmt"
	"github.com/lazyledger/nmt/namespace"
)

func Test_rowRootsFromNamespaceID(t *testing.T) {
	data := generateRandNamespacedRawData(16, 8, 8)
	nID := data[len(data)/2][:8]
	dah, err := makeDAHeader(data)
	if err != nil {
		t.Fatal(err)
	}

	indices, err := rowRootsFromNamespaceID(nID, dah)
	if err != nil {
		t.Fatal(err)
	}

	expected := make([]int, 0)
	for i, row := range dah.RowsRoots {
		if !namespace.ID(nID).Less(row.Min) && namespace.ID(nID).LessOrEqual(row.Max) {
			expected = append(expected, i)
		}
	}

	assert.Equal(t, expected, indices)
}

func Test_unsuccessful_rowRootsFromNamespaceID(t *testing.T) {
	data := generateRandNamespacedRawData(16, 8, 8)
	nIDBelow := make([]byte, 8)
	for i, byt := range data[0][:8] {
		// assuming the data is never 0
		nIDBelow[i] = byt - 1
	}
	nIDExceeds := make([]byte, 8)
	for i, byt := range data[len(data)-1][:8] {
		nIDExceeds[i] = byt + 1
	}

	dah, err := makeDAHeader(data)
	if err != nil {
		t.Fatal(err)
	}

	indicesBelow, err := rowRootsFromNamespaceID(nIDBelow, dah)

	assert.Equal(t, 1, len(indicesBelow))
	assert.Equal(t, 0, indicesBelow[0])
	assert.True(t, strings.Contains(err.Error(), "below minimum"))

	indicesExceeds, err := rowRootsFromNamespaceID(nIDExceeds, dah)

	assert.Equal(t, 1, len(indicesExceeds))
	assert.Equal(t, len(data)-1, indicesExceeds[0])
	require.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "exceeds maximum"))
}

func makeDAHeader(data [][]byte) (*types.DataAvailabilityHeader, error) {
	rows, err := types.NmtRootsFromBytes(data)
	if err != nil {
		return nil, err
	}
	clns, err := types.NmtRootsFromBytes(data)
	if err != nil {
		return nil, err
	}

	return &types.DataAvailabilityHeader{
		RowsRoots:   rows,
		ColumnRoots: clns,
	}, nil
}

func TestRetrieveShares(t *testing.T) {
	api := mockedIpfsAPI(t)

	type test struct {
		name                 string
		expectedErr          error
		nID                  []byte
		data                 types.Data
		startIndex, endIndex int
	}

	tests := []test{
		{
			name: "single transaction",
			nID:  consts.TxNamespaceID,
			data: types.Data{
				Txs: generateRandomContiguousShares(1),
			},
			startIndex: 0,
			endIndex:   1,
		},
		{
			name: "many transactions",
			nID:  consts.TxNamespaceID,
			data: types.Data{
				Txs: generateRandomContiguousShares(20),
			},
			startIndex: 0,
			endIndex:   20,
		},
		{
			name: "single message mixed with contiguous shares",
			nID:  []byte{1, 2, 3, 4, 5, 6, 7, 8},
			data: types.Data{
				Txs: generateRandomContiguousShares(2),
				Messages: types.Messages{
					MessagesList: []types.Message{
						{
							NamespaceID: []byte{1, 2, 3, 4, 5, 6, 7, 8},
							Data:        bytes.Repeat([]byte{1}, consts.MsgShareSize-2),
						},
					},
				},
			},
			startIndex: 2,
			endIndex:   3,
		},
		{
			name: "multiple messages that have the same namespace",
			nID:  []byte{1, 2, 3, 4, 5, 6, 7, 8},
			data: types.Data{
				Txs: generateRandomContiguousShares(2),
				Messages: types.Messages{
					MessagesList: []types.Message{
						{
							NamespaceID: []byte{1, 2, 3, 4, 5, 6, 7, 8},
							Data:        bytes.Repeat([]byte{1}, 300),
						},
						{
							NamespaceID: []byte{1, 2, 3, 4, 5, 6, 7, 8},
							Data:        bytes.Repeat([]byte{1}, 100),
						},
						{
							NamespaceID: []byte{1, 2, 3, 4, 5, 6, 7, 8},
							Data:        bytes.Repeat([]byte{1}, 42),
						},
					},
				},
			},
			startIndex: 2,
			endIndex:   6,
		},
		{
			name: "find between message that crosses multiple shares",
			nID:  []byte{1, 2, 3, 4, 5, 6, 7, 8},
			data: types.Data{
				Txs: generateRandomContiguousShares(2),
				Messages: types.Messages{
					MessagesList: []types.Message{
						{
							NamespaceID: []byte{1, 2, 3, 4, 5, 6, 7, 7},
							Data:        bytes.Repeat([]byte{1}, 300),
						},
						{
							NamespaceID: []byte{1, 2, 3, 4, 5, 6, 7, 8},
							Data:        bytes.Repeat([]byte{1}, 600),
						},
						{
							NamespaceID: []byte{1, 2, 3, 4, 5, 6, 7, 9},
							Data:        bytes.Repeat([]byte{1}, 42),
						},
					},
				},
			},
			startIndex: 4,
			endIndex:   7,
		},
		{
			name:        "missing namespace not found",
			expectedErr: ErrNotFoundInRange,
			nID:         []byte{1, 2, 3, 4, 5, 6, 7, 8},
			data: types.Data{
				Txs: generateRandomContiguousShares(2),
				Messages: types.Messages{
					MessagesList: []types.Message{
						{
							NamespaceID: []byte{1, 2, 3, 4, 5, 6, 7, 7},
							Data:        bytes.Repeat([]byte{1}, 300),
						},
						{
							NamespaceID: []byte{1, 2, 3, 4, 5, 6, 7, 9},
							Data:        bytes.Repeat([]byte{1}, 42),
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			tc := tc
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
			defer cancel()

			block := types.Block{
				LastCommit: &types.Commit{},
				Data:       tc.data,
			}

			// fill the data availability header
			block.Hash()

			// save the block to the dag
			err := PutBlock(ctx, api.Dag(), &block, ipfs.MockRouting(), log.TestingLogger())
			require.NoError(t, err)

			// retrieve desired shares
			result, err := RetrieveShares(ctx, tc.nID, &block.DataAvailabilityHeader, api)
			if tc.expectedErr != nil {
				require.Equal(t, tc.expectedErr, err)
				return
			}
			require.NoError(t, err)

			// compare with original
			shares, _ := tc.data.ComputeShares()
			rawShares := shares.RawShares()
			require.Equal(t, rawShares[tc.startIndex:tc.endIndex], result)
		})
	}

}

// todo fix later
func commitTreeDataToDAG(ctx context.Context, data [][]byte, batchAdder *NmtNodeAdder) (namespace.IntervalDigest, error) {

	tree := nmt.New(sha256.New, nmt.NodeVisitor(batchAdder.Visit)) // TODO consider changing this to default size
	// add some fake data
	for _, d := range data {
		if err := tree.Push(d); err != nil {
			panic(fmt.Sprintf("unexpected error: %v", err))
		}
	}
	return tree.Root(), nil
}

func Test_multipleLeaves_findStartingIndex(t *testing.T) {
	// set nID
	api := mockedIpfsAPI(t)
	ctx := context.Background()
	treeRoots := make(types.NmtRoots, 0)

	var (
		leaves [][]byte
		nID    []byte
	)

	// create nmt adder wrapping batch adder
	batchAdder := NewNmtNodeAdder(ctx, format.NewBatch(ctx, api.Dag()))

	for i := 0; i < 4; i++ {
		data := generateRandNamespacedRawData(4, nmt.DefaultNamespaceIDLen, 16)
		if len(nID) == 0 {
			index := rand.Intn(len(data) - 2)
			leaves = make([][]byte, 2)
			for i := range leaves {
				leaves[i] = make([]byte, 24)
			}
			nID = make([]byte, 8)
			// TODO make this nicer later
			copy(nID, data[index][:8])
			// make 2 byte slices in data have same nID so that nID associated with multiple leaves
			if index == len(data)-2 {
				copy(data[index-1], append(nID, data[index-1][8:]...))

				copy(leaves[0], data[index-1])
				copy(leaves[1], data[index])
				fmt.Println("data index: ", data[index], "\n data index-1: ", data[index-1])
				fmt.Println("leaf 0: ", leaves[0], "\nleaf 1: ", leaves[1])
			} else {
				copy(data[index+1], append(nID, data[index+1][8:]...))
				copy(leaves[0], data[index])
				copy(leaves[1], data[index+1])
				fmt.Println("data index: ", data[index], "\n data index+1: ", data[index+1])
			}
		}
		fmt.Printf("%+v\n", data)

		treeRoot, err := commitTreeDataToDAG(ctx, data, batchAdder)
		if err != nil {
			t.Fatal(err)
		}
		treeRoots = append(treeRoots, treeRoot)

		if err := batchAdder.Commit(); err != nil {
			t.Fatal(err)
		}
	}
	fmt.Println("NID: ", nID)

	dah := &types.DataAvailabilityHeader{
		RowsRoots: treeRoots,
	}

	rootCid, err := plugin.CidFromNamespacedSha256(dah.RowsRoots[0].Bytes())
	if err != nil {
		t.Error(err)
	}
	fmt.Println("rootCID: ", rootCid)
	shares, err := getSharesByNamespace(ctx, nID, dah, []int{0}, api)
	if err != nil {
		t.Fatal(err)
	}
	for i, share := range shares {
		assert.Equal(t, leaves[i], share, i)
	}
}

func Test_successful_findStartingIndex(t *testing.T) {
	// set nID
	api := mockedIpfsAPI(t)
	treeRoots := make(types.NmtRoots, 0)

	ctx := context.Background()

	var (
		nIDData []byte
		nID     []byte
	)

	// create nmt adder wrapping batch adder
	batchAdder := NewNmtNodeAdder(ctx, format.NewBatch(ctx, api.Dag()))

	for i := 0; i < 4; i++ {
		data := generateRandNamespacedRawData(4, nmt.DefaultNamespaceIDLen, 16)
		fmt.Printf("%+v\n", data)
		if len(nID) == 0 {
			nIDData = data[rand.Intn(len(data)-1)] // todo maybe make this nicer later
			nID = nIDData[:8]
		}

		treeRoot, err := commitTreeDataToDAG(ctx, data, batchAdder)
		if err != nil {
			t.Fatal(err)
		}
		treeRoots = append(treeRoots, treeRoot)

		if err := batchAdder.Commit(); err != nil {
			t.Fatal(err)
		}
	}

	fmt.Println("NID DATA: ", nIDData, "\nnID: ", nID)

	dah := &types.DataAvailabilityHeader{
		RowsRoots: treeRoots,
	}

	rootCid, err := plugin.CidFromNamespacedSha256(dah.RowsRoots[0].Bytes())
	if err != nil {
		t.Error(err)
	}
	fmt.Println("rootCID: ", rootCid)

	startingIndex, err := findStartingIndex(ctx, nID, dah, rootCid, api)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("starting index: ", startingIndex)

	leaf, err := GetLeafData(ctx, rootCid, uint32(startingIndex), uint32(len(dah.RowsRoots)), api.Dag())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, nIDData, leaf)
}

func Test_unsuccessful_findStartingIndex(t *testing.T) {
	// set nID
	api := mockedIpfsAPI(t)
	treeRoots := make(types.NmtRoots, 0)

	ctx := context.Background()

	var (
		nIDData []byte
		nID     []byte
	)

	// create nmt adder wrapping batch adder
	batchAdder := NewNmtNodeAdder(ctx, format.NewBatch(ctx, api.Dag()))

	for i := 0; i < 4; i++ {
		data := generateRandNamespacedRawData(4, nmt.DefaultNamespaceIDLen, 16)
		fmt.Printf("%+v\n", data)
		if len(nID) == 0 {
			nIDData = data[rand.Intn(len(data)-1)] // todo maybe make this nicer later
			nID = nIDData[:8]
		}

		treeRoot, err := commitTreeDataToDAG(ctx, data, batchAdder)
		if err != nil {
			t.Fatal(err)
		}
		treeRoots = append(treeRoots, treeRoot)

		if err := batchAdder.Commit(); err != nil {
			t.Fatal(err)
		}
	}

	fmt.Println("NID DATA: ", nIDData, "\nnID: ", nID)

	dah := &types.DataAvailabilityHeader{
		RowsRoots: treeRoots,
	}

	// get rootCID of a row in which nID does NOT exist
	rootCid, err := plugin.CidFromNamespacedSha256(dah.RowsRoots[2].Bytes())
	if err != nil {
		t.Error(err)
	}
	fmt.Println("rootCID: ", rootCid)

	_, err = findStartingIndex(ctx, nID, dah, rootCid, api)
	if err == nil {
		t.Fatal(err)
	}
	fmt.Println(err.Error())
	assert.True(t, strings.Contains(err.Error(), "within range of namespace IDs in tree, but does not exist"))
}

func Test_startIndexFromPath(t *testing.T) {
	var tests = []struct {
		path     []string
		expected int
	}{
		{
			path:     []string{"0", "0", "1"},
			expected: 1,
		},
		{
			path:     []string{"0", "1", "1", "1"},
			expected: 7,
		},
		{
			path:     []string{"0", "0"},
			expected: 0,
		},
		{
			path:     []string{"1", "1", "0"},
			expected: 6,
		},
		{
			path:     []string{"0", "1", "0", "1"},
			expected: 5,
		},
		{
			path:     []string{"0", "0", "0", "0"},
			expected: 0,
		},
		{
			path:     []string{"1", "1", "1", "1"},
			expected: 15,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := startIndexFromPath(tt.path)
			if got != tt.expected {
				t.Fatalf("expected %d, got %d", tt.expected, got)
			}
		})
	}
}

func mockedIpfsAPI(t *testing.T) iface.CoreAPI {
	node, err := ipfs.Mock()()
	if err != nil {
		panic(err)
	}

	ipfsAPI, err := coreapi.NewCoreAPI(node)
	if err != nil {
		panic(err)
	}

	return ipfsAPI
}
