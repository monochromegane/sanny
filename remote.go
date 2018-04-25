package sanny

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"net"
)

type Remote struct {
	address  string
	compress bool
	conn     net.Conn
}

func NewRemote(address string, compress bool) *Remote {
	return &Remote{
		address:  address,
		compress: compress,
	}
}

func (r *Remote) Build(data [][]float32) {
	conn, err := r.getConn()
	if err != nil {
		panic(err)
	}
	for i, _ := range data {
		conn.Write(Float32sToBytes(data[i]))
		buf := make([]byte, SIZE_FLOAT32)
		_, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}
		if ret := int(int32(binary.BigEndian.Uint32(buf))); ret != 0 {
			panic(ret)
		}
	}
	r.Close()
	r.getConn() // Workaround for query warm-up
}

func (r *Remote) Search(q []float32, n int) []int {
	conn, err := r.getConn()
	if err != nil {
		panic(err)
	}

	req := Float32sToBytes(q)
	if r.compress {
		buf := new(bytes.Buffer)
		w := zlib.NewWriter(buf)
		w.Write(req)
		w.Close()
		req = buf.Bytes()
	}

	conn.Write(req)
	buf := make([]byte, SIZE_INT32*n)
	len, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}
	buf = buf[:len]
	return BytesToInts(buf)
}

func (r *Remote) Close() {
	if r.conn == nil {
		return
	}
	r.conn.Close()
	r.conn = nil
}

func (r *Remote) getConn() (net.Conn, error) {
	if r.conn == nil {
		conn, err := net.Dial("tcp", r.address)
		if err != nil {
			return nil, err
		}
		r.conn = conn
	}
	return r.conn, nil
}
