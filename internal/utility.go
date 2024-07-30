package internal

import (
	"flag"
	"log"
)

func GetCommandLineInputs() (i, o, usr, pwd, key, kh string, p, t, w int) {
	iFile := flag.String("f", "input.json", "input file name")
	oFile := flag.String("o", "output.json", "output file name")
	noWorkers := flag.Int("w", 4, "number of workers")
	port := flag.Int("p", 22, "generic default ssh port")
	timeout := flag.Int("t", 30, "timeout [secs]")
	user := flag.String("usr", "admin", "generic username for connection")
	passwd := flag.String("passwd", "admin", "generic password for connection")
	key_path := flag.String("rsa", "", ".ssh file path: ($HOME)/.ssh/id_rsa")
	hosts_path := flag.String("knownhosts", "", ".ssh known hosts file path:($HOME)/.ssh/known_hosts")

	flag.Parse()

	log.Println("File accepted:", *iFile, "| output file:", *oFile)
	log.Println("workers:", *noWorkers, "| timeout:", *timeout, "| port:", *port, "| username:", *user, "| password:", *passwd)
	log.Println("key path:", *key_path, "| knowhosts:", *hosts_path)

	return *iFile, *oFile, *user, *passwd, *key_path, *hosts_path, *port, *timeout, *noWorkers
}

func Execute(input Input, usr, passwd, rsa, kh string, p, t int) (out Output) {
	// binding input if empty from json
	if input.Port == 0 {
		input.Port = p
	}
	if input.Timeout == 0 {
		input.Timeout = t
	}
	if input.Username == "" {
		input.Username = usr
	}
	if input.Password == "" {
		input.Password = passwd
	}

	// doing the ssh to end device
	result, errors := doSSH(input, rsa, kh)

	// forming result Output
	out.I = input
	out.Err = errors
	out.O = result
	return
}
