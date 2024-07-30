package internal

import (
	"context"
	"flag"
	"log"
	"sync"
	"time"
)

func GetCommandLineInputs() CmdPipe {
	iFile := flag.String("f", "input.json", "input file name")
	oFile := flag.String("o", "output.json", "output file name")
	noWorkers := flag.Int("w", 4, "number of workers")
	port := flag.Int("p", 22, "generic default ssh port")
	timeout := flag.Int("t", 30, "timeout [secs]")
	circleTimeout := flag.Int("ct", 120, "timeout [secs]")
	user := flag.String("usr", "admin", "generic username for connection")
	passwd := flag.String("passwd", "admin", "generic password for connection")
	key_path := flag.String("rsa", "", ".ssh file path: ($HOME)/.ssh/id_rsa")
	hosts_path := flag.String("knownhosts", "", ".ssh known hosts file path:($HOME)/.ssh/known_hosts")

	flag.Parse()

	log.Println("File accepted:", *iFile, "| output file:", *oFile)
	log.Println("workers:", *noWorkers, "| port:", *port, "| username:", *user, "| password:", *passwd)
	log.Println("timeout:", *timeout, "| circle timeout:", *circleTimeout, "| key path:", *key_path, "| knowhosts:", *hosts_path)

	return CmdPipe{
		Ifile: *iFile, Ofile: *oFile, Workers: *noWorkers, Port: *port,
		Timeout: *timeout, CircleTimeout: *circleTimeout, Usr: *user, Passwd: *passwd,
		KeyPath: *key_path, HostsPath: *hosts_path,
	}
}

func run(input Input, cmdAgrs CmdPipe, ind int) (out Output) {
	// binding input if empty from json
	if input.Port == 0 {
		input.Port = cmdAgrs.Port
	}
	if input.Timeout == 0 {
		input.Timeout = cmdAgrs.Timeout
	}
	if input.Username == "" {
		input.Username = cmdAgrs.Usr
	}
	if input.Password == "" {
		input.Password = cmdAgrs.Passwd
	}

	// doing the ssh to end device
	result, errors := doSSH(ind, input, cmdAgrs.KeyPath, cmdAgrs.HostsPath)

	// forming result Output
	out.I = input
	out.Err = errors
	out.O = result
	return
}

func Execute(inputs []Input, ch chan<- Output, cmdAgrs CmdPipe) (err error) {
	var wg sync.WaitGroup
	c := make(chan int, cmdAgrs.Workers)
	for i := 0; i < len(inputs); i++ {
		wg.Add(1)
		c <- 1
		go func(input Input, cmdAgrs CmdPipe, ind int) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cmdAgrs.CircleTimeout))
			defer func() { wg.Done(); cancel(); <-c }()

			chInner := make(chan Output)
			go func() {
				defer close(chInner)
				_output := run(input, cmdAgrs, ind)
				chInner <- _output
			}()

			select {
			case _output := <-chInner:
				ch <- _output
			case <-ctx.Done():
				// close(chInner)
				var _output = Output{I: input, Err: []string{ctx.Err().Error()}}
				ch <- _output
				log.Println(ind, input.Host, "Terminate Error", ctx.Err())
			}

		}(inputs[i], cmdAgrs, i+1)
		log.Println("Host sent for ssh -- ", inputs[i], i+1)
	}
	wg.Wait()
	return nil
}
