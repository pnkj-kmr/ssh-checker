package main

import (
	"log"
	// _ "net/http/pprof"
	"ssh-checker/internal"
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
	cmdArgs := internal.GetCommandLineInputs()

	// init
	checker := internal.NewJSON(cmdArgs.Ifile, cmdArgs.Ofile)
	inputs := checker.GetInput()
	log.Println("Total items ----- ", len(inputs))

	// Setting up the core channels
	exitCh := make(chan struct{})
	ch := make(chan internal.Output, cmdArgs.Workers)
	go checker.ProduceOutput(ch, exitCh)

	// Calling main execution to proceed
	err := internal.Execute(inputs, ch, cmdArgs)
	if err != nil {
		log.Fatal("Error during main execution -", err)
	}
	close(ch)
	<-exitCh
	log.Println("Execution completed. time taken", time.Since(st))
}
