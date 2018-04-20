package sanny

import "sort"

type KeyDistance struct {
	sort.Interface
	idx []int
}

func (kd KeyDistance) Swap(i, j int) {
	kd.Interface.Swap(i, j)
	kd.idx[i], kd.idx[j] = kd.idx[j], kd.idx[i]
}

func NewKeyDistance(n sort.Interface) *KeyDistance {
	kd := &KeyDistance{Interface: n, idx: make([]int, n.Len())}
	for i := range kd.idx {
		kd.idx[i] = i
	}
	return kd
}

func NewIntKeyDistance(n ...int) *KeyDistance         { return NewKeyDistance(sort.IntSlice(n)) }
func NewFloat64KeyDistance(n ...float64) *KeyDistance { return NewKeyDistance(sort.Float64Slice(n)) }
