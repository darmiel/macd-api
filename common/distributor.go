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
