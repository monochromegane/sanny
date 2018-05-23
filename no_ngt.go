// +build !ngt

package sanny

type NGT struct {
}

func NewNGT(edge int) *NGT {
	return &NGT{}
}

func (ngt *NGT) Build(data [][]float32) {
	panic("Not support searcher. Please build with `ngt` tag.")
}

func (ngt NGT) Search(q []float32, n int) []int {
	panic("Not support searcher. Please build with `ngt` tag.")
}
