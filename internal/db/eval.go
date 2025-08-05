package db

import (
	"strings"

	par "github.com/razzat008/letsgodb/internal/Parser"
)

// EvalWhere evaluates a WHERE expression (Expr) against a row.
// columns: schema column names
// row: row values (as []string)
func EvalWhere(expr par.Expr, columns, row []string) bool {
	switch e := expr.(type) {
	case *par.Condition:
		// Find column index
		idx := -1
		for i, col := range columns {
			if col == e.Column {
				idx = i
				break
			}
		}
		if idx == -1 || idx >= len(row) {
			return false
		}
		val := strings.Trim(row[idx], "'")
		condVal := strings.Trim(e.Value, "'")
		switch e.Operator {
		case "=":
			return val == condVal
		case "!=":
			return val != condVal
		case ">":
			return val > condVal
		case "<":
			return val < condVal
		case ">=":
			return val >= condVal
		case "<=":
			return val <= condVal
		default:
			return false
		}
	case *par.BinaryExpr:
		left := EvalWhere(e.Left, columns, row)
		right := EvalWhere(e.Right, columns, row)
		switch strings.ToUpper(e.Operator) {
		case "AND":
			return left && right
		case "OR":
			return left || right
		default:
			return false
		}
	default:
		return false
	}
}
