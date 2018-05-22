package sanny

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
