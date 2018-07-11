package kafka

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {

	hasher := NewAsc()
	for i := 1; i < 40; i++ {
		str := fmt.Sprintf("w%vv.plat.bjtb.pdtv.it", i)
		hasher.Reset()
		hasher.Write([]byte(str))
		t.Logf("[%v] hash result %v", i, hasher.Sum32()%20)
	}
}
