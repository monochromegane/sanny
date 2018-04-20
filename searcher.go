package sanny

type Searcher interface {
	Build([][]float32)
	Search([]float32, int) []int
}
