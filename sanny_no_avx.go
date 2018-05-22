// +build !avx

package sanny

func bruteForceAlgorism() Searcher {
	return NewBruteForceBLAS()
}
