package db

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"

	"github.com/razzat008/letsgodb/internal/storage"
)

// serializing row : converts a slice of string value into a length prefixed byte slice
// format: [row_length(uint16)][csv_data]
func SerializeRow(values []string) []byte {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	_ = writer.Write(values)
	writer.Flush()
	rowBytes := buf.Bytes()
	length := uint16(len(rowBytes))
	out := make([]byte, 2+len(rowBytes))
	binary.LittleEndian.PutUint16(out[0:2], length)
	copy(out[2:], rowBytes)
	return out
}

// deserializeRow reads a length-prefixed row from data and returns the values and bytes consumed.
func DeserializeRow(data []byte) ([]string, int) {
	if len(data) < 2 {
		return nil, 0
	}
	length := binary.LittleEndian.Uint16(data[0:2])
	if len(data) < int(2+length) {
		return nil, 0
	}
	rowBytes := data[2 : 2+length]
	reader := csv.NewReader(bytes.NewReader(rowBytes))
	values, err := reader.Read()
	if err != nil {
		return nil, 0
	}
	return values, int(2 + length)
}

// InsertRow appends a row to the last page with space, or allocates a new page if needed.
// Returns the page number where the row was written.
func InsertRow(pager *storage.Pager, values []string) (uint32, error) {
	rowBytes := SerializeRow(values)
	// For now, always use page 0 (expand later for multi-page)
	var pageNum uint32 = 0
	if pager.PageCount() == 0 {
		pageNum = pager.AllocatePage()
	}
	page := pager.GetPage(pageNum)

	// Find offset: scan for first zero byte (simple, not robust for deletes)
	offset := 0
	for offset < storage.PageSize {
		if page[offset] == 0 {
			break
		}
		// Skip to next row
		_, consumed := DeserializeRow(page[offset:])
		if consumed == 0 {
			break
		}
		offset += consumed
	}

	if offset+len(rowBytes) > storage.PageSize {
		// Not enough space, allocate new page
		pageNum = pager.AllocatePage()
		page = pager.GetPage(pageNum)
		offset = 0
	}
	copy(page[offset:], rowBytes)
	err := pager.FlushPage(pageNum, page)
	return pageNum, err
}

// ReadAllRows reads all rows from all pages in the pager.
func ReadAllRows(pager *storage.Pager) [][]string {
	var rows [][]string
	for pageNum := uint32(0); pageNum < uint32(pager.PageCount()); pageNum++ {
		page := pager.GetPage(pageNum)
		offset := 0
		for offset < storage.PageSize {
			values, consumed := DeserializeRow(page[offset:])
			if consumed == 0 || values == nil || (len(values) > 0 && values[0] == "") {
				break
			}
			rows = append(rows, values)
			offset += consumed
		}
	}
	return rows
}
