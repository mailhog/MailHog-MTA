package smtp

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type fakeRw struct {
	_read  func(p []byte) (n int, err error)
	_write func(p []byte) (n int, err error)
	_close func() error
}

func (rw *fakeRw) Read(p []byte) (n int, err error) {
	if rw._read != nil {
		return rw._read(p)
	}
	return 0, nil
}
func (rw *fakeRw) Close() error {
	if rw._close != nil {
		return rw._close()
	}
	return nil
}
func (rw *fakeRw) Write(p []byte) (n int, err error) {
	if rw._write != nil {
		return rw._write(p)
	}
	return len(p), nil
}

func TestAccept(t *testing.T) {
	Convey("Accept should handle a connection", t, func() {
		frw := &fakeRw{}
		Accept("1.1.1.1:11111", frw, "localhost")
	})
}

func TestSocketError(t *testing.T) {
	Convey("Socket errors should return from Accept", t, func() {
		frw := &fakeRw{
			_read: func(p []byte) (n int, err error) {
				return -1, errors.New("OINK")
			},
		}
		Accept("1.1.1.1:11111", frw, "localhost")
	})
}
