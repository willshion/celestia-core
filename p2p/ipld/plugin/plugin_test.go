package lazyledger_ipld

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	shell "github.com/ipfs/go-ipfs-api"
	config "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	icore "github.com/ipfs/interface-go-ipfs-core"

	"github.com/lazyledger/nmt"
	"github.com/lazyledger/rsmt2d"
)

func TestDataSquareRowOrColumnRawInputParserCidEqNmtRoot(t *testing.T) {
	ctx := context.Background()
	api := spawnEphemeral(ctx, t)
	dagService := api.Dag()

	tests := []struct {
		name     string
		leafData [][]byte
	}{
		{"16 leaves", generateRandNamespacedRawData(16, namespaceSize, shareSize)},
		{"32 leaves", generateRandNamespacedRawData(32, namespaceSize, shareSize)},
		{"extended row", generateExtendedRow(t)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			n := nmt.New(sha256.New())
			buf := createByteBufFromRawData(t, tt.leafData)
			for _, share := range tt.leafData {
				err := n.Push(share[:namespaceSize], share[namespaceSize:])
				if err != nil {
					t.Errorf("nmt.Push() unexpected error = %v", err)
					return
				}
			}
			gotNodes, err := DataSquareRowOrColumnRawInputParser(buf, 0, 0)
			if err != nil {
				t.Errorf("DataSquareRowOrColumnRawInputParser() unexpected error = %v", err)
				return
			}
			rootNodeCid := gotNodes[0].Cid()
			multiHashOverhead := 2
			lastNodeHash := rootNodeCid.Hash()
			if got, want := lastNodeHash[multiHashOverhead:], n.Root().Bytes(); !bytes.Equal(got, want) {
				t.Errorf("hashes don't match\ngot: %v\nwant: %v", got, want)
			}
			lastNodeCid := gotNodes[len(gotNodes)-1].Cid()
			if gotHash, wantHash := lastNodeCid.Hash(), hashLeaf(tt.leafData[0]); !bytes.Equal(gotHash[multiHashOverhead:], wantHash) {
				t.Errorf("first node's hash does not match the Cid\ngot: %v\nwant: %v", gotHash[multiHashOverhead:], wantHash)
			}
			nodePrefixOffset := 1 // leaf / inner node prefix is one byte
			lastLeafNodeData := gotNodes[len(gotNodes)-1].RawData()
			if gotData, wantData := lastLeafNodeData[nodePrefixOffset:], tt.leafData[0]; !bytes.Equal(gotData, wantData) {
				t.Errorf("first node's data does not match the leaf's data\ngot: %v\nwant: %v", gotData, wantData)
			}

			err = dagService.AddMany(ctx, gotNodes)
			if err != nil {
				t.Errorf("APIDagService.AddMany(): %v", err)
			}
			gotNode, err := dagService.Get(ctx, rootNodeCid)
			if err != nil {
				t.Errorf("dagService.Get(%v): %v", rootNodeCid, err)
			}
			root := gotNodes[0]
			if !reflect.DeepEqual(gotNode, root) {
				t.Errorf("got diff, got: %v, want: %v", gotNode.(nmtNode), root.(nmtNode))
			}
		})
	}
}

func TestDagPutWithPlugin(t *testing.T) {
	t.Skip("Requires running ipfs daemon with the plugin compiled and installed")

	t.Log("Warning: running this test writes to your local IPFS block store!")

	const numLeaves = 32
	data := generateRandNamespacedRawData(numLeaves, namespaceSize, shareSize)
	buf := createByteBufFromRawData(t, data)
	printFirst := 10
	t.Logf("first leaf, nid: %x, data: %x...", data[0][:namespaceSize], data[0][namespaceSize:namespaceSize+printFirst])
	n := nmt.New(sha256.New())
	for _, share := range data {
		err := n.Push(share[:namespaceSize], share[namespaceSize:])
		if err != nil {
			t.Errorf("nmt.Push() unexpected error = %v", err)
			return
		}
	}

	sh := shell.NewLocalShell()
	cid, err := sh.DagPut(buf, "raw", DagParserFormatName)
	if err != nil {
		t.Fatalf("DagPut() failed: %v", err)
	}
	// convert NMT tree root to CID and verify it matches the CID returned by DagPut
	treeRootBytes := n.Root().Bytes()
	nmtCid, err := cidFromNamespacedSha256(treeRootBytes)
	if err != nil {
		t.Fatalf("cidFromNamespacedSha256() failed: %v", err)
	}
	if nmtCid.String() != cid {
		t.Errorf("CIDs from NMT and plugin do not match: got %v, want: %v", cid, nmtCid.String())
	}
	// print out cid s.t. it can be used on the commandline
	t.Logf("Stored with cid: %v\n", cid)

	// DagGet leaf by leaf:
	for i, wantShare := range data {
		gotLeaf := &nmtLeafNode{}
		path := leafIdxToPath(cid, i)
		if err := sh.DagGet(path, gotLeaf); err != nil {
			t.Errorf("DagGet(%s) failed: %v", path, err)
		}
		if gotShare := gotLeaf.Data; !bytes.Equal(gotShare, wantShare) {
			t.Errorf("DagGet returned different data than pushed, got: %v, want: %v", gotShare, wantShare)
		}
	}
}

func createTempRepo() (string, error) {
	repoPath, err := ioutil.TempDir("", "ipfs-shell")
	if err != nil {
		return "", fmt.Errorf("failed to get temp dir: %s", err)
	}

	// Create a config with default options and a 2048 bit key
	cfg, err := config.Init(ioutil.Discard, 2048)
	if err != nil {
		return "", err
	}

	// Create the repo with the config
	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to init ephemeral node: %s", err)
	}

	return repoPath, nil
}

// Spawns a node to be used just for this run (i.e. creates a tmp repo)
func spawnEphemeral(ctx context.Context, t *testing.T) icore.CoreAPI {
	// Load any external plugins if available on externalPluginsPath
	plugins, err := loader.NewPluginLoader(filepath.Join("", "plugins"))
	if err != nil {
		t.Fatalf("error loading plugins: %s", err)
	}

	// Load preloaded and external plugins
	if err := plugins.Initialize(); err != nil {
		t.Fatalf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		t.Fatalf("error initializing plugins: %s", err)
	}

	// Create a Temporary Repo
	repoPath, err := createTempRepo()
	if err != nil {
		t.Fatalf("failed to create temp repo: %s", err)
	}

	// Spawning an ephemeral IPFS node
	return createNode(ctx, repoPath, t)
}

// Creates an IPFS node and returns its coreAPI
func createNode(ctx context.Context, repoPath string, t *testing.T) icore.CoreAPI {
	// Open the repo
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		t.Fatalf("fsrepo.Open(%v): %v", repoPath, err)
	}

	// Construct the node

	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
		// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		Repo: repo,
	}

	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		t.Fatalf("core.NewNode(): %v", err)
	}

	// Attach the Core API to the constructed node
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		t.Fatalf("coreapi.NewCoreAPI(): %v", err)
	}
	return api
}

func generateExtendedRow(t *testing.T) [][]byte {
	origData := generateRandNamespacedRawData(16, namespaceSize, shareSize)
	origDataWithoutNamespaces := make([][]byte, 16)
	for i, share := range origData {
		origDataWithoutNamespaces[i] = share[namespaceSize:]
	}

	extendedData, err := rsmt2d.ComputeExtendedDataSquare(origDataWithoutNamespaces, rsmt2d.RSGF8, newNmtConstructor)
	if err != nil {
		t.Fatalf("rsmt2d.Encode(): %v", err)
		return nil
	}
	extendedRow := extendedData.Row(0)
	for i, rowCell := range extendedRow {
		if i < len(origData)/4 {
			nid := origData[i][:namespaceSize]
			extendedRow[i] = append(nid, rowCell...)
		} else {
			maxNid := bytes.Repeat([]byte{0xFF}, namespaceSize)
			extendedRow[i] = append(maxNid, rowCell...)
		}
	}
	return extendedRow
}

var _ rsmt2d.Tree = &nmtWrapper{}

func newNmtConstructor() rsmt2d.Tree {
	return &nmtWrapper{
		nmt.New(sha256.New()),
	}
}

// we could get rid of this wrapper and use the nmt directly if we
// make Push take in the data as one one byte array (instead of two).
type nmtWrapper struct {
	*nmt.NamespacedMerkleTree
}

func (n nmtWrapper) Push(data []byte) {
	if err := n.NamespacedMerkleTree.Push(data[:namespaceSize], data[namespaceSize:]); err != nil {
		panic(err)
	}
}

func (n nmtWrapper) Prove(idx int) (merkleRoot []byte, proofSet [][]byte, proofIndex uint64, numLeaves uint64) {
	proof, err := n.NamespacedMerkleTree.Prove(idx)
	if err != nil {
		panic(err)
	}
	return n.NamespacedMerkleTree.Root().Bytes(),
		proof.Nodes(),
		uint64(proof.Start()),
		0 // TODO: NMT doesn't return the number of leaves
}

func (n nmtWrapper) Root() []byte {
	return n.NamespacedMerkleTree.Root().Bytes()
}

func leafIdxToPath(cid string, idx int) string {
	// currently this fmt directive assumes 32 leaves:
	bin := fmt.Sprintf("%05b", idx)
	path := strings.Join(strings.Split(bin, ""), "/")
	return cid + "/" + path
}

func createByteBufFromRawData(t *testing.T, leafData [][]byte) *bytes.Buffer {
	buf := bytes.NewBuffer(make([]byte, 0))
	for _, share := range leafData {
		_, err := buf.Write(share)
		if err != nil {
			t.Fatalf("buf.Write() unexpected error = %v", err)
			return nil
		}
	}
	return buf
}

// this snippet of the nmt internals is copied here:
func hashLeaf(data []byte) []byte {
	h := sha256.New()
	nID := data[:namespaceSize]
	toCommittToDataWithoutNID := data[namespaceSize:]

	res := append(append(make([]byte, 0), nID...), nID...)
	data = append([]byte{nmt.LeafPrefix}, toCommittToDataWithoutNID...)
	h.Write(data)
	return h.Sum(res)
}

func generateRandNamespacedRawData(total int, nidSize int, leafSize int) [][]byte {
	data := make([][]byte, total)
	for i := 0; i < total; i++ {
		nid := make([]byte, nidSize)
		rand.Read(nid)
		data[i] = nid
	}
	sortByteArrays(data)
	for i := 0; i < total; i++ {
		d := make([]byte, leafSize)
		rand.Read(d)
		data[i] = append(data[i], d...)
	}

	return data
}

func sortByteArrays(src [][]byte) {
	sort.Slice(src, func(i, j int) bool { return bytes.Compare(src[i], src[j]) < 0 })
}
