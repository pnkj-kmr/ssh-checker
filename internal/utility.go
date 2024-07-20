package internal

import (
	"flag"
	"fmt"
	"log"
)

func GetCommandLineInputs() (i, o, usr, pwd, key, kh string, p, t, w int) {
	iFile := flag.String("f", "input.json", "input file name")
	oFile := flag.String("o", "output.json", "output file name")
	noWorkers := flag.Int("w", 4, "number of workers")
	port := flag.Int("p", 22, "default ssh port")
	timeout := flag.Int("t", 30, "timeout [secs]")
	user := flag.String("usr", "admin", "username for connection")
	passwd := flag.String("passwd", "admin", "password for connection")
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

	// calling main ssh executor func
	outputChannel := make(chan []byte)
	err := doSSH(input, rsa, kh, outputChannel)
	if err == nil {
		out.Ok = true
		err = fmt.Errorf("")
	}
	out.Err = err.Error()

	// binding result to output
	var result []string
	var str string
	for x := range outputChannel {
		str = string(x)
		result = append(result, str)
	}
	out.I = input
	out.O = result
	return
}

/**

// calling main execution function
		outputChannel := make(chan []byte)
		err = controllers.Execute(payload.Src, payload.Dst, payload.Proto, outputChannel)
		if err != nil {
			errMsg := fmt.Sprintf("Error: %s", err.Error())
			utils.L.Error("Error", zap.String("err", errMsg))
			// skiping error if occurs such cases
			skipError := "remote command exited without exit status or exit signal"
			if !strings.Contains(errMsg, skipError) {
				_ = writeToWS(ws, mt, errMsg, false, true)
			}
		}
		// writing into ws
		var str string
		for x := range outputChannel {
			str = string(x)
			ok, str := checkUnwantedData(str)
			if ok {
				_ = writeToWS(ws, mt, maskingIP(str, payload.Dst, payload.Proto), false, false)
			}
		}
		_ = writeToWS(ws, mt, "", true, false)


*/
