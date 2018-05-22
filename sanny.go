package sanny

import (
	"math"
	"sync"
)

type Sanny struct {
	data      [][]float32
	splitNum  int
	Top       int
	searchers []Searcher
	sort      bool
	indecies  [][]int
}

func NewSanny(splitNum, top int, sort bool, searchers []Searcher, indecies [][]int) Searcher {
	return &Sanny{
		splitNum:  splitNum,
		Top:       top,
		searchers: searchers,
		sort:      sort,
		indecies:  indecies,
	}
}

func (s *Sanny) Build(data [][]float32) {
	if s.sort {
		s.data = data
	}
	dim := len(data[0])
	colNum := int(math.Ceil(float64(dim) / float64(s.splitNum)))
	for i, ids := range s.indecies {
		if _, ok := s.searchers[i].(*Remote); ok {
			// Build on remote server
			continue
		}

		from := ids[0] * colNum
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

func (s Sanny) Search(q []float32, n int) []int {
	results := map[int]int{}

	var wg sync.WaitGroup
	ch := make(chan int, s.Top)
	done := make(chan struct{}, 1)
	go func() {
		for id := range ch {
			results[id] += 1
		}
		done <- struct{}{}
	}()

	dim := len(q)
	colNum := int(math.Ceil(float64(dim) / float64(s.splitNum)))
	for i, ids := range s.indecies {
		wg.Add(1)
		go func(i int, ids []int) {
			defer wg.Done()
			var qq []float32
			if len(ids) == 1 {
				from := ids[0] * colNum
				to := from + colNum
				if to > dim {
					to = dim
				}
				qq = q[from:to]
			} else {
				for j, _ := range ids {
					from := ids[j] * colNum
					to := from + colNum
					if to > dim {
						to = dim
					}
					qq = append(qq, q[from:to]...)
				}
			}
			r := s.searchers[i].Search(qq, s.Top)
			for _, id := range r {
				ch <- id
			}
		}(i, ids)
	}
	wg.Wait()
	close(ch)
	<-done

	return s.candidates(q, n, results)
}

func (s Sanny) candidates(q []float32, n int, candidates map[int]int) []int {
	if s.sort {
		return s.bruteForce(q, n, candidates)
	}

	ids := make([]int, len(candidates))
	i := 0
	for id := range candidates {
		ids[i] = id
		i++
	}
	return ids
}

func (s Sanny) bruteForce(q []float32, n int, candidates map[int]int) []int {
	data := make([][]float32, len(candidates))
	cnt := 0
	keys := make([]int, len(candidates))
	for k, _ := range candidates {
		data[cnt] = s.data[k]
		keys[cnt] = k
		cnt += 1
	}
	bf := bruteForceAlgorism()
	bf.Build(data)
	indecies := bf.Search(q, n)
	ids := make([]int, len(indecies))
	for i, s := range indecies {
		ids[i] = keys[s]
	}
	return ids
}
