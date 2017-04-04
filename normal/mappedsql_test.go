package normal

import "testing"

func TestROL(t *testing.T) {
	nms := new(NormalMappedSql)

	nms.Id = "xxxx-r"

	if !nms.IsReadOnly() {
		t.Fail()
	}

}

func TestROM(t *testing.T) {
	nms := new(NormalMappedSql)

	nms.Id = "select-afdadsf"

	if !nms.IsReadOnly() {
		t.Fail()
	}

}
func TestRON(t *testing.T) {
	nms := new(NormalMappedSql)

	nms.Id = "insert-xa"

	if nms.IsReadOnly() {
		t.Fail()
	}

}
