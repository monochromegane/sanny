package sanny

import (
	"math"
	"sort"
)

type BruteForce struct {
	data [][]float32
}

func NewBruteForce() *BruteForce {
	return &BruteForce{}
}

func (bf *BruteForce) Build(data [][]float32) {
	bf.data = data
}

func (bf BruteForce) Search(q []float32, n int) []int {
	ids, _ := bf.searchWithDistance(q, n)
	return ids
}

func (bf BruteForce) searchWithDistance(q []float32, n int) ([]int, []float64) {
	distances := make([]float64, len(bf.data))
	for i, d := range bf.data {
		distances[i] = bf.distance(q, d)
	}
	s := NewFloat64KeyDistance(distances...)
	sort.Sort(s)
	return s.idx[:n], distances[:n]
}

func (bf BruteForce) distance(a, b []float32) float64 {
	var norm float64
	for i, v := range a {
		diff := float64(b[i] - v)
		norm = math.Hypot(norm, diff)
	}
	return norm
}
