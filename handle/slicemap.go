package handle

import (
	"github.com/it512/sqlt"
	"github.com/jmoiron/sqlx"
)

type (
	SliceMapRowsHandle struct {
		Data [][]map[string]interface{}
	}
)

func (h *SliceMapRowsHandle) AddResultSet() {
	h.Data = append(h.Data, make([]map[string]interface{}, 0, 10))
}

func (h *SliceMapRowsHandle) HandleRow(r sqlt.RowScanner) {
	idx := len(h.Data) - 1
	m := make(map[string]interface{})
	if e := sqlx.MapScan(r, m); e == nil {
		h.Data[idx] = append(h.Data[idx], m)
	}
}
