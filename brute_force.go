package sanny

import (
	"math"
	"runtime"
	"sort"
	"sync"
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

type KeyValues struct {
	Key      int
	Values   []float32
	Distance float64
}

func (bf BruteForce) SearchConcurrent(q []float32, n int) []int {
	var wg sync.WaitGroup
	cpus := runtime.NumCPU()
	in := make(chan KeyValues, cpus)
	out := make(chan KeyValues, cpus)
	for i := 0; i < cpus; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for kv := range in {
				distance := bf.distance(q, kv.Values)
				kv.Distance = distance
				out <- kv
			}
		}()
	}
	go func() {
		for i, d := range bf.data {
			in <- KeyValues{Key: i, Values: d}
		}
		close(in)
		wg.Wait()
		close(out)
	}()

	ch := make(chan []int)
	go func() {
		distances := make([]float64, n)
		idx := make([]int, n)
		for i, _ := range distances {
			distances[i] = math.MaxFloat64
		}
		for kv := range out {
			if kv.Distance > distances[n-1] {
				continue
			}
			distances = append(distances, kv.Distance)
			idx = append(idx, kv.Key)

			s := NewFloat64KeyDistance(distances...)
			s.idx = idx
			sort.Sort(s)

			distances = distances[:n]
			idx = s.idx[:n]
		}
		ch <- idx
	}()

	return <-ch
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
