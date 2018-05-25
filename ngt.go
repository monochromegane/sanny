// +build ngt

package sanny

import (
	"io/ioutil"
	"runtime"

	"github.com/yahoojapan/gongt"
)

type NGT struct {
	edge int
	idx  *gongt.NGT
}

func NewNGT(edge int) *NGT {
	return &NGT{
		edge: edge, // gongt.DefaultSearchEdgeSize = 40
	}
}

func (ngt *NGT) Build(data [][]float32) {
	dim := len(data[0])
	dir, _ := ioutil.TempDir("", "NGT")
	idx := gongt.New(dir).SetObjectType(gongt.Float).SetDimension(dim).SetSearchEdgeSize(ngt.edge).Open()
	for _, d := range data {
		obj := make([]float64, dim)
		for j, _ := range d {
			obj[j] = float64(d[j])
		}
		idx.Insert(obj)
	}
	idx.CreateIndex(runtime.NumCPU())
	ngt.idx = idx
}

func (ngt NGT) Search(q []float32, n int) []int {
	var result []int
	q64 := make([]float64, len(q))
	for i, _ := range q {
		q64[i] = float64(q[i])
	}
	r, err := ngt.idx.Search(q64, n, gongt.DefaultEpsilon)
	if err != nil {
		return result
	}
	for i, _ := range r {
		result = append(result, r[i].ID-1)
	}
	return result
}
