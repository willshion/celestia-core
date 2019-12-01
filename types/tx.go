package types

import (
	"bytes"
	"errors"
	"fmt"

	amino "github.com/tendermint/go-amino"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
	cmn "github.com/tendermint/tendermint/libs/common"
)

// Tx is an arbitrary byte array.
// NOTE: Tx has no types at this level, so when wire encoded it's just length-prefixed.
// Might we want types here ?
//
// LAZY: TXs also need a namespace in LL.
// While this could be achieved leaving this a blob and letting the
// (abci) application deal with the content of the TXs, we can't really do so
// in LL as the DataHash field in the header needs to be namespaced.
type Tx []byte

//
// Emulate LL messages: Namespace returns the namespace of a message
// (of size namespaceSize).
func (tx Tx) Namespace() []byte {
	return tx[:namespaceSize]
}

// Emulate LL messages: Data returns the data of a message.
func (tx Tx) Data() []byte {
	return tx[namespaceSize:]
}

// Hash computes the TMHASH hash of the wire encoded transaction.
func (tx Tx) Hash() []byte {
	return tmhash.Sum(tx)
}

func (tx Tx) NamespacedString() string {
	return fmt.Sprintf("Tx{namespace: %X, data: %X}", tx.Namespace(), tx.Data())
}

// String returns the hex-encoded transaction as a string.
func (tx Tx) String() string {
	return fmt.Sprintf("Tx{%X}", []byte(tx))
}

// Txs is a slice of Tx.
// LAZY: Txs.Hash() is used to compute that DataHash in the header (merkle root of the list
// messages/Txs).
//
// Note: for a 1st LL prototype with as little changes as possible, we can
// simply assume that a Tx is of the form ([namespaceSize]byte..., []data ...).
// Hence we can always extract the namespace and the actual Tx.
type Txs []Tx

// Hash returns the Merkle root hash of the transaction hashes.
// i.e. the leaves of the tree are the hashes of the txs.
func (txs Txs) Hash() []byte {
	// These allocations will be removed once Txs is switched to [][]byte,
	// ref #2603. This is because golang does not allow type casting slices without unsafe
	txBzs := make([][]byte, len(txs))
	for i := 0; i < len(txs); i++ {
		txBzs[i] = txs[i].Hash()
	}
	// LAZY: For the simplistic validity rules (section 4.1) this just
	// needs to be a "namespaced" merkle tree instead of a "simple merkle tree".
	//
	//
	// from the lazyledger paper:
	//
	// each block header h_i contains the root mRoot_i of a Merkle tree of a
	// list of messages M_i = (m^0_i , m^1_i , ...), such that given
	// a function root(M) that returns the Merkle root of a list of messages M,
	// then root(M_i) = mRoot_i.
	// This is not an ordinary Merkle tree,
	// but an ordered Merkle tree we refer to as a ‘namespaced’ Merkle tree [...]
	//
	//
	// LAZY: For the probability validity rule (section 4.2),
	// we need to compute column and row roots.
	//
	// As a general note, it would be great if tendermint made it easy to swap the
	// underlying merkle tree (or even sth. more general).
	// Additionally, it would be good if tendermint would abstract away from the
	// underlying block structure.
	// Ideally, tendermint would rather just reach consensus on "something" instead of
	// assuming the exact structure of what this something is. This opaque data on which
	// consensus is reached could have validity rules. In the case of a block this could
	// be the simple block validity rule.
	return merkle.SimpleHashFromByteSlices(txBzs)
}

// Index returns the index of this transaction in the list, or -1 if not found
func (txs Txs) Index(tx Tx) int {
	for i := range txs {
		if bytes.Equal(txs[i], tx) {
			return i
		}
	}
	return -1
}

// IndexByHash returns the index of this transaction hash in the list, or -1 if not found
func (txs Txs) IndexByHash(hash []byte) int {
	for i := range txs {
		if bytes.Equal(txs[i].Hash(), hash) {
			return i
		}
	}
	return -1
}

// Proof returns a simple merkle proof for this node.
// Panics if i < 0 or i >= len(txs)
// TODO: optimize this!
func (txs Txs) Proof(i int) TxProof {
	l := len(txs)
	bzs := make([][]byte, l)
	for i := 0; i < l; i++ {
		bzs[i] = txs[i].Hash()
	}
	root, proofs := merkle.SimpleProofsFromByteSlices(bzs)

	return TxProof{
		RootHash: root,
		Data:     txs[i],
		Proof:    *proofs[i],
	}
}

// TxProof represents a Merkle proof of the presence of a transaction in the Merkle tree.
type TxProof struct {
	RootHash cmn.HexBytes       `json:"root_hash"`
	Data     Tx                 `json:"data"`
	Proof    merkle.SimpleProof `json:"proof"`
}

// Leaf returns the hash(tx), which is the leaf in the merkle tree which this proof refers to.
func (tp TxProof) Leaf() []byte {
	return tp.Data.Hash()
}

// Validate verifies the proof. It returns nil if the RootHash matches the dataHash argument,
// and if the proof is internally consistent. Otherwise, it returns a sensible error.
func (tp TxProof) Validate(dataHash []byte) error {
	if !bytes.Equal(dataHash, tp.RootHash) {
		return errors.New("Proof matches different data hash")
	}
	if tp.Proof.Index < 0 {
		return errors.New("Proof index cannot be negative")
	}
	if tp.Proof.Total <= 0 {
		return errors.New("Proof total must be positive")
	}
	valid := tp.Proof.Verify(tp.RootHash, tp.Leaf())
	if valid != nil {
		return errors.New("Proof is not internally consistent")
	}
	return nil
}

// TxResult contains results of executing the transaction.
//
// One usage is indexing transaction results.
type TxResult struct {
	Height int64                  `json:"height"`
	Index  uint32                 `json:"index"`
	Tx     Tx                     `json:"tx"`
	Result abci.ResponseDeliverTx `json:"result"`
}

// ComputeAminoOverhead calculates the overhead for amino encoding a transaction.
// The overhead consists of varint encoding the field number and the wire type
// (= length-delimited = 2), and another varint encoding the length of the
// transaction.
// The field number can be the field number of the particular transaction, or
// the field number of the parenting struct that contains the transactions []Tx
// as a field (this field number is repeated for each contained Tx).
// If some []Tx are encoded directly (without a parenting struct), the default
// fieldNum is also 1 (see BinFieldNum in amino.MarshalBinaryBare).
func ComputeAminoOverhead(tx Tx, fieldNum int) int64 {
	fnum := uint64(fieldNum)
	typ3AndFieldNum := (fnum << 3) | uint64(amino.Typ3_ByteLength)
	return int64(amino.UvarintSize(typ3AndFieldNum)) + int64(amino.UvarintSize(uint64(len(tx))))
}
