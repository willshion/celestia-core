package flaghasher

import (
	"hash"

	"github.com/LazyLedger/lazyledger-prototype"
)

// Mostly
const flagSize = 16

var codedFlag [flagSize]byte

func init() {
	for i := range codedFlag {
		codedFlag[i] = 0xFF
	}
}

type flagDigest struct {
	flagger    lazyledger.Flagger
	baseHasher hash.Hash
	data       []byte
	codedMode  bool
}

// New returns a new hash.Hash computing checksums using the baseHasher with flags from flagger.
func New(flagger lazyledger.Flagger, baseHasher hash.Hash) *flagDigest {
	return &flagDigest{
		flagger:    flagger,
		baseHasher: baseHasher,
	}
}

// TODO: add proper comment here!
func (d *flagDigest) SetCodedMode(mode bool) {
	d.codedMode = mode
}

func (d *flagDigest) Write(p []byte) (int, error) {
	d.data = append(d.data, p...)
	return d.baseHasher.Write(p)
}

func (d *flagDigest) Sum(in []byte) []byte {
	in = append(in, d.parentFlag()...)
	return d.baseHasher.Sum(in)
}

func (d *flagDigest) Size() int {
	return d.flagger.FlagSize() + d.baseHasher.Size()
}

func (d *flagDigest) BlockSize() int {
	return d.baseHasher.BlockSize()
}

func (d *flagDigest) Reset() {
	d.data = nil
	d.baseHasher.Reset()
}

func (d *flagDigest) leftFlag() []byte {
	return d.flagger.NodeFlag(d.data[1 : d.Size()+1])
}

func (d *flagDigest) rightFlag() []byte {
	return d.flagger.NodeFlag(d.data[1+d.Size():])
}

func (d *flagDigest) parentFlag() []byte {
	if d.isLeaf() {
		if d.codedMode {
			return codedFlag[:]
		}
		return d.flagger.LeafFlag(d.mainData())
	}
	return d.flagger.Union(d.leftFlag(), d.rightFlag())
}

func (d *flagDigest) mainData() []byte {
	return d.data[1:]
}

func (d *flagDigest) isLeaf() bool {
	return d.data[0] == byte(0)
}
