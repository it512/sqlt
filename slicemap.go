package sqlt

type (
	SliceMapRowsHandle struct {
		Data         [][]map[string]interface{}
		keyConverter func(string) string
	}
)

func (h *SliceMapRowsHandle) AddResultSet() {
	h.Data = append(h.Data, make([]map[string]interface{}, 0, 10))
}

func (h *SliceMapRowsHandle) HandleRow(r RowScanner) {
	idx := len(h.Data) - 1
	m := make(map[string]interface{})
	if e := scan(r, h.keyConverter, m); e == nil {
		h.Data[idx] = append(h.Data[idx], m)
	}
}

func (h *SliceMapRowsHandle) Count() int {
	return len(h.Data)
}

func (h *SliceMapRowsHandle) ResuleSet(i int) []map[string]interface{} {
	return h.Data[i]
}

func scan(r RowScanner, c func(string) string, dest map[string]interface{}) error {
	columns, err := r.Columns()
	if err != nil {
		return err
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	err = r.Scan(values...)
	if err != nil {
		return err
	}

	for i, column := range columns {
		dest[c(column)] = *(values[i].(*interface{}))
	}

	return r.Err()
}

func NewSliceMapRowsHandler(keyConverter func(string) string) *SliceMapRowsHandle {
	return &SliceMapRowsHandle{Data: make([][]map[string]interface{}, 0), keyConverter: keyConverter}
}
