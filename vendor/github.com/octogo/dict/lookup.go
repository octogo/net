package dict

type lookup struct {
	Key      string
	Response chan interface{}
	Error    chan error
}

func (l *lookup) Close() {
	close(l.Response)
	close(l.Error)
}
