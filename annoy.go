// +build annoy

package sanny

import "annoyindex"

type Annoy struct {
	tree    int
	SearchK int
	idx     annoyindex.AnnoyIndexEuclidean
}

func NewAnnoy(tree, searchK int) *Annoy {
	return &Annoy{
		tree:    tree,
		SearchK: searchK,
	}
}

func (a *Annoy) Build(data [][]float32) {
	idx := annoyindex.NewAnnoyIndexEuclidean(len(data[0]))
	for i, d := range data {
		idx.AddItem(i, d)
	}
	idx.Build(a.tree)
	a.idx = idx
}

func (a Annoy) Search(q []float32, n int) []int {
	var result []int
	a.idx.GetNnsByVector(q, n, a.SearchK, &result)
	return result
}
