package customhttp

import (
	"fmt"
	"testing"
)

type fakeConn struct {
	request []byte //Method Path Data
	input   []byte
	output  []byte
	closed  bool
}

func (f *fakeConn) Read(rbuf []byte) (int, error) {
	fmt.Println("read called on fake conn")
	copy(rbuf[:len(f.request)], f.request)
	f.input = rbuf
	return len(rbuf), nil
}

func (f *fakeConn) Write(wbuf []byte) (int, error) {
	fmt.Println("write called on fake conn")
	fmt.Printf("data written %v", string(wbuf))
	f.output = wbuf
	return len(wbuf), nil
}

func (f *fakeConn) Close() error {
	fmt.Println("close called on fake conn")
	f.closed = true
	return nil
}

func TestHandleRequest_GET_Success(t *testing.T) {
	f := &fakeConn{
		request: []byte("GET /test"),
		closed:  false,
	}

	routes := map[string]string{"GET": "/test"}

	chttp := InitHttp(f, routes)
	chttp.HandleRequest()
	if !f.closed {
		t.Log("expected connection to be closed")
		t.Fail()
	}

	if len(f.input) == 0 {
		t.Log("expected input bytes")
		t.Fail()
	}
	fmt.Printf("received req: %v", string(f.input))

	if len(f.output) == 0 {
		t.Log("expected output bytes")
		t.Fail()
	}
	fmt.Printf("outgoin res: %v", string(f.output))
}
