package main

import (
	"fmt"

	"github.com/monochromegane/sanny"
)

func main() {
	data := [][]float32{
		[]float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0},
		[]float32{10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0, 80.0, 90.0},
		[]float32{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9},
	}

	// Server
	build := make(chan struct{})
	go func() {
		sdata := make([][]float32, len(data))
		for i, _ := range data {
			sdata[i] = data[i][3:]
		}
		searcher := sanny.NewSanny(
			2,
			2,
			false,
			[]sanny.Searcher{
				sanny.NewAnnoy(1, 10),
				sanny.NewAnnoy(1, 10),
			},
			[][]int{[]int{0}, []int{1}},
		)
		searcher.Build(sdata)
		server := sanny.NewServer("localhost:8000", false, searcher)
		build <- struct{}{}
		server.Run()
	}()

	// Client
	searcher := sanny.NewSanny(
		3,
		2,
		true,
		[]sanny.Searcher{
			sanny.NewAnnoy(1, 10),
			sanny.NewRemote("localhost:8000", false),
		},
		[][]int{[]int{0}, []int{1, 2}},
	)

	searcher.Build(data)

	<-build
	ids := searcher.Search([]float32{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9}, 2)
	fmt.Printf("%v\n", ids)
}
