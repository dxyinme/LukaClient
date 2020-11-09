package IpcMsg

import (
	"testing"
)

func TestArrayBufferToByteArray(t *testing.T) {
	var (
		test string
		expect []byte
		result []byte
	)
	checker := func(A,B *[]byte) bool {
		if len(*A) != len(*B) {
			return false
		}
		for i := 0 ; i < len(*A); i ++ {
			if (*A)[i] != (*B)[i] {
				return false
			}
		}
		return true
	}
	test = "1,200,93,73,4,158,58,5,65,242,65,128,55,135,217,129,191,128,0,16"
	expect = []byte{1,200,93,73,4,158,58,5,65,242,65,128,55,135,217,129,191,128,0,16}
	result = ArrayBufferToByteArray(&test)
	if !checker(&result, &expect) {
		t.Errorf("test %d failed.", 1)
	}
}