package main

import (
	"log"
	// _ "net/http/pprof"
	"ssh-checker/internal"
	"sync"
	"time"
)

// func cpu_profile() {
// 	f, err := os.Create("cpuprofile")
// 	if err != nil {
// 		log.Fatal("could not create CPU profile: ", err)
// 	}
// 	defer f.Close() // error handling omitted for example
// 	if err := pprof.StartCPUProfile(f); err != nil {
// 		log.Fatal("could not start CPU profile: ", err)
// 	}
// 	defer pprof.StopCPUProfile()
// }

// func mem_profile() {
// 	f, err := os.Create("memprofile")
// 	if err != nil {
// 		log.Fatal("could not create memory profile: ", err)
// 	}
// 	defer f.Close() // error handling omitted for example
// 	runtime.GC()    // get up-to-date statistics
// 	if err := pprof.WriteHeapProfile(f); err != nil {
// 		log.Fatal("could not write memory profile: ", err)
// 	}
// }

func main() {
	// add for debug purpose
	// cpu_profile()
	// mem_profile()

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
		go func(input internal.Input, usr, passwd, rsa, kh string, p, t, ind int) {
			defer func() { wg.Done(); <-c }()
			out := internal.Execute(input, usr, passwd, rsa, kh, p, t, ind)
			ch <- out
		}(inputs[i], usr, passwd, rsa, kh, p, t, i+1)
		log.Println("Host sent for ssh -- ", inputs[i], i+1)
	}
	wg.Wait()
	close(ch)
	<-exitCh
	log.Println("Execution completed. time taken", time.Since(st))
}
