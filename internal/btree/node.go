package btree

import (
	"bytes"
	"encoding/binary"
)

const (
	PageSize         = 4096
	NodeTypeLeaf     = 0
	NodeTypeInternal = 1

	// Offsets in the page
	OffsetNodeType  = 0  // 1 byte
	OffsetKeyCount  = 1  // 2 bytes (uint16)
	OffsetParent    = 3  // 4 bytes (uint32)
	OffsetNextLeaf  = 7  // 4 bytes (uint32, only for leaf)
	OffsetDataStart = 11 // start of keys/values
)

// NodeHeader represents the metadata at the start of a page.
type NodeHeader struct {
	NodeType byte
	KeyCount uint16
	Parent   uint32
	NextLeaf uint32 // only for leaf nodes
}

// ReadNodeHeader reads the header from a page.
func ReadNodeHeader(page []byte) NodeHeader {
	return NodeHeader{
		NodeType: page[OffsetNodeType],
		KeyCount: binary.LittleEndian.Uint16(page[OffsetKeyCount : OffsetKeyCount+2]),
		Parent:   binary.LittleEndian.Uint32(page[OffsetParent : OffsetParent+4]),
		NextLeaf: binary.LittleEndian.Uint32(page[OffsetNextLeaf : OffsetNextLeaf+4]),
	}
}

// WriteNodeHeader writes the header to a page.
func WriteNodeHeader(page []byte, h NodeHeader) {
	page[OffsetNodeType] = h.NodeType
	binary.LittleEndian.PutUint16(page[OffsetKeyCount:OffsetKeyCount+2], h.KeyCount)
	binary.LittleEndian.PutUint32(page[OffsetParent:OffsetParent+4], h.Parent)
	binary.LittleEndian.PutUint32(page[OffsetNextLeaf:OffsetNextLeaf+4], h.NextLeaf)
}

// InitLeafNode initializes a page as a leaf node.
func InitLeafNode(page []byte) {
	WriteNodeHeader(page, NodeHeader{
		NodeType: NodeTypeLeaf,
		KeyCount: 0,
		Parent:   0,
		NextLeaf: 0,
	})
}

// InitInternalNode initializes a page as an internal node.
func InitInternalNode(page []byte) {
	WriteNodeHeader(page, NodeHeader{
		NodeType: NodeTypeInternal,
		KeyCount: 0,
		Parent:   0,
		NextLeaf: 0, // unused for internal
	})
}

// InsertKeyValueLeaf inserts a key/value pair into a leaf node page.
// Returns false if not enough space.
func InsertKeyValueLeaf(page []byte, key, value []byte) bool {
	header := ReadNodeHeader(page)
	pos := OffsetDataStart

	// Find insert position (sorted by key)
	for i := 0; i < int(header.KeyCount); i++ {
		k, _, next := readKeyValueAt(page, pos)
		if compareKeys(key, k) < 0 {
			break
		}
		pos = next
	}

	// Calculate space needed
	needed := 2 + len(key) + 2 + len(value)
	if len(page)-freeSpaceLeaf(page) < needed {
		return false // Not enough space
	}

	// Shift existing data to make room
	end := OffsetDataStart + usedSpaceLeaf(page)
	copy(page[pos+needed:end+needed], page[pos:end])

	// Write new key/value
	writeKeyValueAt(page, pos, key, value)

	// Update header
	header.KeyCount++
	WriteNodeHeader(page, header)
	return true
}

// Helper: read key/value at offset, return key, value, and next offset
func readKeyValueAt(page []byte, offset int) ([]byte, []byte, int) {
	keyLen := int(binary.LittleEndian.Uint16(page[offset : offset+2]))
	key := page[offset+2 : offset+2+keyLen]
	valOffset := offset + 2 + keyLen
	valLen := int(binary.LittleEndian.Uint16(page[valOffset : valOffset+2]))
	val := page[valOffset+2 : valOffset+2+valLen]
	next := valOffset + 2 + valLen
	return key, val, next
}

// Helper: write key/value at offset
func writeKeyValueAt(page []byte, offset int, key, value []byte) {
	binary.LittleEndian.PutUint16(page[offset:offset+2], uint16(len(key)))
	copy(page[offset+2:offset+2+len(key)], key)
	valOffset := offset + 2 + len(key)
	binary.LittleEndian.PutUint16(page[valOffset:valOffset+2], uint16(len(value)))
	copy(page[valOffset+2:valOffset+2+len(value)], value)
}

// Helper: compare keys (lexicographically)
func compareKeys(a, b []byte) int {
	return bytes.Compare(a, b)
}

// Helper: calculate used space in leaf node
func usedSpaceLeaf(page []byte) int {
	header := ReadNodeHeader(page)
	pos := OffsetDataStart
	for i := 0; i < int(header.KeyCount); i++ {
		_, _, next := readKeyValueAt(page, pos)
		pos = next
	}
	return pos - OffsetDataStart
}

// Helper: calculate free space in leaf node
func freeSpaceLeaf(page []byte) int {
	return PageSize - OffsetDataStart - usedSpaceLeaf(page)
}
