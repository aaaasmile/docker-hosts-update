// Code generated by "stringer -type TokenType"; DO NOT EDIT.

package info

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[itemCommentKey-0]
	_ = x[itemText-1]
	_ = x[itemIP-2]
	_ = x[itemName-3]
	_ = x[itemError-4]
	_ = x[itemEOF-5]
}

const _TokenType_name = "itemCommentKeyitemTextitemIPitemNameitemErroritemEOF"

var _TokenType_index = [...]uint8{0, 14, 22, 28, 36, 45, 52}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
