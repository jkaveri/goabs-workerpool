package main

import (
	"fmt"
	"sync"
)

func main() {
	a := make(chan struct{}, 1)
	var wg sync.WaitGroup
	go func() {
		<- a
		fmt.Println("closed")
		wg.Done()
		return
	}()
	wg.Add(1);
	close(a)

	wg.Wait()

}
