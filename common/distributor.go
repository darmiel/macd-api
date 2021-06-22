package common

import "sync"

func DistributedGoroutine(data []interface{}, num int, exec func([]interface{})) {
	var wg sync.WaitGroup
	var current []interface{}
	exe := func(c []interface{}) {
		wg.Add(1)
		go func() {
			exec(c)
			wg.Done()
		}()
	}
	for _, d := range data {
		current = append(current, d)
		if len(current) >= num {
			exe(current)
			current = nil
		}
	}
	if len(current) > 0 {
		exe(current)
		current = nil
	}
	wg.Wait()
}

// 3, 2, 1, 2, 3, 4, 2, 3, 4, 3, 2, 1
// 0, 1, 2, 3     4, 5, 6, 7     8, 9
func Slice(size, step int) (resp [][]int) {
	var idx []int
	for i := 0; i < size; i++ {
		idx = append(idx, i)
		if len(idx) >= step {
			resp = append(resp, idx)
			idx = nil
		}
	}
	if len(idx) != 0 {
		resp = append(resp, idx)
	}
	return
}
