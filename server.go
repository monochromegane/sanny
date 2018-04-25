package sanny

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
)

type Server struct {
	compress bool
	port     int
	searcher Searcher
	dim      int
}

func NewServer(port, dim int, compress bool, searcher Searcher) Server {
	return Server{
		port:     port,
		dim:      dim,
		compress: compress,
		searcher: searcher,
	}
}

func (s *Server) Initialize() error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}
	defer ln.Close()

	conn, err := ln.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()

	var data [][]float32
	for {
		buf := make([]byte, 4*s.dim)
		len, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		buf = buf[:len]
		data = append(data, BytesToFloat32s(buf))
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes[0:4], uint32(0))
		conn.Write(bytes)
	}
	s.searcher.Build(data)
	return nil
}

func (s Server) Run() error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go func(conn net.Conn) {
			for {
				buf := make([]byte, 4*s.dim)
				len, err := conn.Read(buf)
				if err == io.EOF {
					break
				}
				buf = buf[:len]
				if s.compress {
					r, _ := zlib.NewReader(bytes.NewReader(buf))
					dc, _ := ioutil.ReadAll(r)
					buf = dc
				}
				q := BytesToFloat32s(buf)
				ids := s.searcher.Search(q, 10)
				conn.Write(IntsToBytes(ids))
			}
			conn.Close()
		}(conn)
	}
}
