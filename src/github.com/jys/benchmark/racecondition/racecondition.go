package racecondition

import "sync"

var foo = 0

func Race() {
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 50000; i++ {
			foo++
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 50000; i++ {
			foo++
		}
	}()
	wg.Wait()
	println(foo)
}
