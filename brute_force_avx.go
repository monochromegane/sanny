// +build avx

package sanny

import (
	"math"
	"runtime"
	"sort"
	"sync"
)

func (bf BruteForce) searchWithDistance(q []float32, n int) ([]int, []float64) {
	distances := make([]float64, len(bf.data))

	dim := len(q)
	x := mmMalloc(dim)
	y := mmMalloc(dim)
	z := mmMalloc(dim)
	defer mmFree(x)
	defer mmFree(y)
	defer mmFree(z)
	for i, _ := range q {
		x[i] = q[i]
	}

	for i, d := range bf.data {
		for j, _ := range d {
			y[j] = d[j]
		}
		avxSub(dim, x, y, z)
		distances[i] = math.Sqrt(float64(avxDot(dim, z, z)))
	}
	s := NewFloat64KeyDistance(distances...)
	sort.Sort(s)
	return s.idx[:n], distances[:n]
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

			dim := len(q)
			x := mmMalloc(dim)
			y := mmMalloc(dim)
			z := mmMalloc(dim)
			defer mmFree(x)
			defer mmFree(y)
			defer mmFree(z)
			for i, _ := range q {
				x[i] = q[i]
			}
			for kv := range in {
				for j, _ := range kv.Values {
					y[j] = kv.Values[j]
				}
				avxSub(dim, x, y, z)
				distance := math.Sqrt(float64(avxDot(dim, z, z)))
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
