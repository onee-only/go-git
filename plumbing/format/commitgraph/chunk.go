package commitgraph

import "bytes"

const (
	szChunkSig = 4 // Length of a chunk signature
)

// chunkSignatures contains the coalesced byte signatures for each chunk type.
// The order of the signatures must match the order of the ChunkType constants.
// (When adding new chunk types you must avoid introducing ambiguity, and you may need to add padding separators to this list or reorder these signatures.)
// (i.e. it would not be possible to add a new chunk type with the signature "IDFO" without some reordering or the addition of separators.)
var chunkSignatures = [...]byte{
	'O', 'I', 'D', 'F',
	'O', 'I', 'D', 'L',
	'C', 'D', 'A', 'T',
	'G', 'D', 'A', '2',
	'G', 'D', 'O', '2',
	'E', 'D', 'G', 'E',
	'B', 'I', 'D', 'X',
	'B', 'D', 'A', 'T',
	'B', 'A', 'S', 'E',
	'\x00', '\x00', '\x00', '\x00',
}

// ChunkType represents the type of a chunk in the commit graph file.
type ChunkType int

// Chunk types in the commit graph file.
const (
	OIDFanoutChunk              ChunkType = iota // "OIDF"
	OIDLookupChunk                               // "OIDL"
	CommitDataChunk                              // "CDAT"
	GenerationDataChunk                          // "GDA2"
	GenerationDataOverflowChunk                  // "GDO2"
	ExtraEdgeListChunk                           // "EDGE"
	BloomFilterIndexChunk                        // "BIDX"
	BloomFilterDataChunk                         // "BDAT"
	BaseGraphsListChunk                          // "BASE"
	ZeroChunk                                    // "\000\000\000\000"
)
const lenChunks = int(ZeroChunk) // ZeroChunk is not a valid chunk type, but it is used to determine the length of the chunk type list.

// Signature returns the byte signature for the chunk type as array.
func (ct ChunkType) Signature() [szChunkSig]byte {
	if ct >= BaseGraphsListChunk || ct < 0 { // not a valid chunk type just return ZeroChunk
		ct = ZeroChunk
	}

	offset := ct * szChunkSig

	var sig [szSignature]byte
	copy(sig[:], chunkSignatures[offset:offset+szChunkSig])

	return sig
}

// SignatureBytes returns the byte signature for the chunk type as slice.
func (ct ChunkType) SignatureBytes() []byte {
	sig := ct.Signature()
	return sig[:]
}

// ChunkTypeFromBytes returns the chunk type for the given byte signature.
func ChunkTypeFromSig(b [szChunkSig]byte) (ChunkType, bool) {
	for idx := 0; idx < len(chunkSignatures); idx += szChunkSig {
		if bytes.Equal(b[:], chunkSignatures[idx:idx+szChunkSig]) {
			return ChunkType(idx / szChunkSig), true
		}
	}

	return -1, false
}
