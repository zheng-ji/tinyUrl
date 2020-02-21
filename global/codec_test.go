package global

import (
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	rawInt := int64(43)
	encodeStr := Encode(rawInt)
	decodeInt := Decode(encodeStr)

	if decodeInt != rawInt {
		t.Error("encode decode fail")
	}
}
