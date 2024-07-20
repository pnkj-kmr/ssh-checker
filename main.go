package main

import (
	"log"
	"ssh-checker/internal"
	"sync"
	"time"
)

func main() {
	st := time.Now()
	i, o, usr, passwd, rsa, kh, p, t, w := internal.GetCommandLineInputs()

	checker := internal.NewJSON(i, o)
	inputs := checker.GetInput()
	log.Println("Total Items ----- ", len(inputs))

	exitCh := make(chan struct{})
	ch := make(chan internal.Output, w)
	go checker.ProduceOutput(ch, exitCh)

	var wg sync.WaitGroup
	c := make(chan int, w)
	for i := 0; i < len(inputs); i++ {
		wg.Add(1)
		c <- 1
		go func(input internal.Input, usr, passwd, rsa, kh string, p, t int) {
			defer func() { wg.Done(); <-c }()
			out := internal.Execute(input, usr, passwd, rsa, kh, p, t)
			ch <- out
		}(inputs[i], usr, passwd, rsa, kh, p, t)
		log.Println("Host sent for ssh -- ", inputs[i], i+1)
	}
	wg.Wait()
	close(ch)
	<-exitCh
	log.Println("Execution completed. time taken", time.Since(st))
}
