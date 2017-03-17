package sqlt

func oops(e error) {
	if e != nil {
		panic(e)
	}
}

func MustExec(i int64, e error) int64 {
	oops(e)
	return i
}

func MustTx(tx *TxOp, e error) *TxOp {
	oops(e)
	return tx
}
