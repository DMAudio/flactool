package util

import "p20190417/types"

func SSPairListSort(input [][2]string, priority func(string) int, reverse bool) [][2]string {
	priorityList := make(types.SNList, len(input))
	for i := 0; i < len(input); i++ {
		priorityList[i].Key = i
		priorityList[i].Priority = priority(input[i][0])
	}

	priorityListSorted := priorityList.Sort(reverse)

	result := make([][2]string, len(input))
	for i, pItem := range *priorityListSorted {
		index := pItem.Key.(int)
		result[i][0] = input[index][0]
		result[i][1] = input[index][1]
	}

	return result
}
