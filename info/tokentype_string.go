// Code generated by "stringer -type TokenType"; DO NOT EDIT.

package info

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[itemCommentKey-0]
	_ = x[itemComment-1]
	_ = x[itemText-2]
	_ = x[itemName-3]
	_ = x[itemIPPart-4]
	_ = x[itemDotKey-5]
	_ = x[itemSpaceSepa-6]
	_ = x[itemTrail-7]
	_ = x[itemError-8]
	_ = x[itemEOF-9]
}

const _TokenType_name = "itemCommentKeyitemCommentitemTextitemNameitemIPPartitemDotKeyitemSpaceSepaitemTrailitemErroritemEOF"

var _TokenType_index = [...]uint8{0, 14, 25, 33, 41, 51, 61, 74, 83, 92, 99}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
