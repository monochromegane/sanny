package sanny

import (
	"bytes"
	"compress/zlib"
	"io"
	"io/ioutil"

	"github.com/monochromegane/smux"
)

type Server struct {
	compress bool
	searcher Searcher
	address  string
}

func NewServer(address string, compress bool, searcher Searcher) Server {
	return Server{
		address:  address,
		compress: compress,
		searcher: searcher,
	}
}

func (s Server) Run() error {
	server := smux.Server{
		Network: "tcp",
		Address: s.address,
		Handler: smux.HandlerFunc(func(w io.Writer, r io.Reader) {
			buf, _ := ioutil.ReadAll(r)
			if s.compress {
				zr, _ := zlib.NewReader(bytes.NewReader(buf))
				dc, _ := ioutil.ReadAll(zr)
				buf = dc
			}
			q := BytesToFloat32s(buf)
			ids := s.searcher.Search(q, 10)
			w.Write(IntsToBytes(ids))
		}),
	}
	return server.ListenAndServe()
}
