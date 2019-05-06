package types

import "sort"

type SNList []struct {
	Key      interface{}
	Priority int
}

func (p SNList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p SNList) Len() int           { return len(p) }
func (p SNList) Less(i, j int) bool { return p[i].Priority < p[j].Priority }
func (p SNList) Sort(Reverse bool) *SNList {
	if Reverse {
		sort.Sort(sort.Reverse(p))
	} else {
		sort.Sort(p)
	}
	return &p
}
