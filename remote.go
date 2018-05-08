package sanny

import (
	"bytes"
	"compress/zlib"

	"github.com/monochromegane/smux"
)

type Remote struct {
	address  string
	compress bool
	client   smux.Client
}

func NewRemote(address string, compress bool) *Remote {
	return &Remote{
		client:   smux.Client{Network: "tcp", Address: address},
		address:  address,
		compress: compress,
	}
}

func (r *Remote) Build(data [][]float32) {
}

func (r *Remote) Search(q []float32, n int) []int {
	req := Float32sToBytes(q)
	if r.compress {
		buf := new(bytes.Buffer)
		w := zlib.NewWriter(buf)
		w.Write(req)
		w.Close()
		req = buf.Bytes()
	}

	res, _ := r.client.Post(req)
	return BytesToInts(res)
}
