package btree

import (
	"errors"

	"github.com/razzat008/letsgodb/internal/storage"
)

// BTree represents a B-Tree for a table.
type BTree struct {
	Pager    *storage.Pager
	RootPage uint32
}

// NewBTree creates a new BTree instance given a pager and root page number.
func NewBTree(pager *storage.Pager, rootPage uint32) *BTree {
	return &BTree{
		Pager:    pager,
		RootPage: rootPage,
	}
}

// Insert inserts a key/value pair into the B-Tree.
// For now, only supports single-leaf (no splits, no internal nodes).
func (bt *BTree) Insert(key, value []byte) error {
	page := bt.Pager.GetPage(bt.RootPage)
	header := ReadNodeHeader(page)
	if header.NodeType != NodeTypeLeaf {
		return errors.New("root is not a leaf node (internal nodes not yet supported)")
	}
	ok := InsertKeyValueLeaf(page, key, value)
	if !ok {
		return errors.New("leaf node full (splitting not yet implemented)")
	}
	return bt.Pager.FlushPage(bt.RootPage, page)
}

// Search looks for a key in the B-Tree and returns its value if found.
// For now, only supports single-leaf (no internal nodes).
func (bt *BTree) Search(key []byte) ([]byte, error) {
	page := bt.Pager.GetPage(bt.RootPage)
	header := ReadNodeHeader(page)
	if header.NodeType != NodeTypeLeaf {
		return nil, errors.New("root is not a leaf node (internal nodes not yet supported)")
	}
	pos := OffsetDataStart
	for i := 0; i < int(header.KeyCount); i++ {
		k, v, next := readKeyValueAt(page, pos)
		if compareKeys(key, k) == 0 {
			return v, nil
		}
		pos = next
	}
	return nil, errors.New("key not found")
}
