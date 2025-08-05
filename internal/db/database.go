package db

import (
	"encoding/binary",
	"strings"
)

import "github.com/razzat008/letsgodb/internal/storage"

// serializing row : converts a slice of string value into a length prefixed byte slice
// format: [row_length(uint16)][csv_data]
func SerializeRow(values []string) []byte {
	row:= strings.Join(values, ",")
	rowBytes := []byte(row) // converting into a byte slice
	length := uint16(len(rowBytes)) // length of the row bytes
	buf:=make([]byte, 2+len(rowBytes))
	binary.LittleEndian.PutUint16(buf[0:2], length)
	copy(buf[2:], rowBytes)
	return buf
}

// deserializeRow reads a length-prefixed row from data and returns the values and bytes consumed.
func DeserializeRow(data []byte) ([]string, int) {
    if len(data) < 2 { // need to have at least 2 bytes for length
        return nil, 0
    }
    length := binary.LittleEndian.Uint16(data[0:2])
    if len(data) < int(2+length) { // need to have at least 2 bytes for length and length bytes for data
        return nil, 0
    }
    rowBytes := data[2 : 2+length]
    row := string(rowBytes)
    values := strings.Split(row, ",")
    return values, int(2 + length)
}
