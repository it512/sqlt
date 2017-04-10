package handler

import (
	"github.com/it512/sqlt"
	"github.com/jmoiron/sqlx"
)

type (
	SliceMapRowHandler struct {
		Data [][]map[string]interface{}
	}
)

func (h *SliceMapRowHandler) AddResultSet() {
	h.Data = append(h.Data, make([]map[string]interface{}, 0, 10))
}

func (h *SliceMapRowHandler) HandleRow(r sqlt.RowScanner) {
	idx := len(h.Data) - 1
	m := make(map[string]interface{})
	if e := sqlx.MapScan(r, m); e == nil {
		h.Data[idx] = append(h.Data[idx], m)
	}
}
