package sanny

import (
	"strconv"

	"github.com/ekzhu/lsh"
)

type LSH struct {
	l   int
	m   int
	w   float64
	idx *lsh.LshForest
}

func NewLSH(l, m int, w float64) *LSH {
	return &LSH{l: l, m: m, w: w}
}

func (l *LSH) Build(data [][]float32) {
	dim := len(data[0])
	idx := lsh.NewLshForest(dim, l.l, l.m, l.w)
	for i, _ := range data {
		obj := make(lsh.Point, dim)
		for j, _ := range data[i] {
			obj[j] = float64(data[i][j])
		}
		idx.Insert(obj, strconv.Itoa(i))
	}
	l.idx = idx
}

func (l LSH) Search(q []float32, n int) []int {
	var result []int
	q64 := make(lsh.Point, len(q))
	for i, _ := range q {
		q64[i] = float64(q[i])
	}
	r := l.idx.Query(q64, n)
	for i, _ := range r {
		id, _ := strconv.Atoi(r[i])
		result = append(result, id)
	}
	return result
}
