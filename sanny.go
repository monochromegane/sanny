package sanny

import (
	"math"
	"sync"
)

type Sanny struct {
	data      [][]float32
	splitNum  int
	top       int
	searchers []Searcher
}

func NewSanny(splitNum, top int, searchers []Searcher) Sanny {
	return Sanny{
		splitNum:  splitNum,
		top:       top,
		searchers: searchers,
	}
}

func (s *Sanny) Build(data [][]float32) {
	s.data = data
	dim := len(data[0])
	colNum := int(math.Ceil(float64(dim) / float64(s.splitNum)))
	for i := 0; i < s.splitNum; i++ {
		from := i * colNum
		to := from + colNum
		if to > dim {
			to = dim
		}
		d := make([][]float32, len(data))
		for j, _ := range data {
			d[j] = data[j][from:to]
		}
		s.searchers[i].Build(d)
	}
}

func (s Sanny) Search(q []float32, n int, distance bool) ([]int, []float64) {
	results := map[int]int{}

	var wg sync.WaitGroup
	ch := make(chan int, s.top)
	done := make(chan struct{}, 1)
	go func() {
		for id := range ch {
			results[id] += 1
		}
		done <- struct{}{}
	}()

	dim := len(q)
	colNum := int(math.Ceil(float64(dim) / float64(s.splitNum)))
	for i := 0; i < s.splitNum; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			from := i * colNum
			to := from + colNum
			if to > dim {
				to = dim
			}
			ids := s.searchers[i].Search(q[from:to], s.top)
			for _, id := range ids {
				ch <- id
			}
		}(i)
	}
	wg.Wait()
	close(ch)
	<-done

	return s.bruteForce(q, n, distance, results)
}

func (s Sanny) bruteForce(q []float32, n int, distance bool, candidates map[int]int) ([]int, []float64) {
	data := make([][]float32, len(candidates))
	cnt := 0
	keys := make([]int, len(candidates))
	for k, _ := range candidates {
		data[cnt] = s.data[k]
		keys[cnt] = k
		cnt += 1
	}
	bf := BruteForceBLAS{}
	bf.Build(data)
	indecies := bf.Search(q, n)
	ids := make([]int, len(indecies))
	for i, s := range indecies {
		ids[i] = keys[s]
	}
	if distance {
		bf := BruteForce{}
		distances := make([]float64, len(ids))
		for i, id := range ids {
			distances[i] = bf.distance(q, s.data[id])
		}
		return ids, distances
	} else {
		return ids, nil
	}
}
