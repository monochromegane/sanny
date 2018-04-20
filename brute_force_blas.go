package sanny

import "sort"

type BruteForceBLAS struct {
	index   [][]float32
	lengths []float32
}

func NewBruteForceBLAS() *BruteForceBLAS {
	return &BruteForceBLAS{}
}

func (bf *BruteForceBLAS) Build(data [][]float32) {
	bf.index = data
	lengths := make([]float32, len(data))
	for i, d := range data {
		var len float32
		for _, c := range d {
			len += c * c
		}
		lengths[i] = len
	}
	bf.lengths = lengths
}

func (bf BruteForceBLAS) Search(q []float32, n int) []int {
	// argmin_a (a - b)^2 = argmin_a a^2 - 2ab + b^2 = argmin_a a^2 - 2ab
	distances := make([]float64, len(bf.index))
	for i, idx := range bf.index {
		distances[i] = float64(bf.lengths[i] - 2*bf.dot(q, idx))
	}
	s := NewFloat64KeyDistance(distances...)
	sort.Sort(s)
	return s.idx[:n]
}

func (bf BruteForceBLAS) dot(a, b []float32) float32 {
	var d float32
	for i, a := range a {
		d += a * b[i]
	}
	return d
}
