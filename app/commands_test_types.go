package main

type DummyConn struct {
	Data []byte
}

func (d *DummyConn) Read(p []byte) (n int, err error) {
	copySize := copy(d.Data, p)
	return copySize, nil
}

func (d *DummyConn) Write(p []byte) (n int, err error) {
	d.Data = append(d.Data, p...)
	return len(d.Data), nil
}

func (d *DummyConn) Close() error {
	return nil
}

func NewTestingContext() *Context {
	return &Context{
		Conn:      &DummyConn{},
		ConnState: NewConnState(),
		State:     NewState(),
	}
}
