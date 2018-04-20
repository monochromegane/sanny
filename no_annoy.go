// +build !annoy

package sanny

type Annoy struct {
}

func NewAnnoy(tree, searchK int) *Annoy {
	return &Annoy{}
}

func (a *Annoy) Build(data [][]float32) {
	panic("Not support searcher. Please build with `annoy` tag.")
}

func (a Annoy) Search(q []float32, n int) []int {
	panic("Not support searcher. Please build with `annoy` tag.")
}
