package internal

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
)

func doSSH(input Input, rsa, kh string, out chan<- []byte) (err error) {
	routineStarted := false
	defer func() {
		if !routineStarted {
			close(out)
		}
	}()

	// TODO - need to pass key and hosts file
	client, err := dialConnection(input, rsa, kh)
	if err != nil {
		log.Println("Connection Error", input.Host, err.Error())
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Println("Session Error", input.Host, err.Error())
		return err
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Println("Session outpipe", input.Host, err.Error())
		return err
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		log.Println("Session errpipe", input.Host, err.Error())
		return err
	}

	// calling listeners
	outChan := make(chan []byte)
	errChan := make(chan []byte)
	outputListener(stdout, stderr, outChan, errChan, out)
	routineStarted = true

	for _, cmd := range input.Commands {
		if err := session.Run(cmd); err != nil {
			log.Println("Command Error", input.Host, err.Error())
			return err
		}
	}

	return nil
}

func dialConnection(input Input, f, kh string) (client *ssh.Client, err error) {
	homeDir := os.Getenv("HOME")
	if f == "" {
		if homeDir == "" {
			homeDir = "/root"
		}
		f = filepath.Join(homeDir, ".ssh", "id_rsa")
	}

	// Reading the ssh rsa file
	key, err := os.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key (%s) : %v", f, err.Error())
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key (%s) : %v", key, err)
	}

	config := &ssh.ClientConfig{
		User:    input.Username,
		Timeout: time.Second * time.Duration(input.Timeout),
		Config: ssh.Config{
			KeyExchanges: preferredKexAlgos,
			Ciphers:      preferredCiphers,
			MACs:         supportedMACs,
		},
		Auth: []ssh.AuthMethod{
			ssh.Password(input.Password), // ssh enable password need to be added
			ssh.PublicKeys(signer),
		},
		//
		// TODO - host key based if needed
		//
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
		// HostKeyCallback: ssh.HostKeyCallback(func(host string, remote net.Addr, pubKey ssh.PublicKey) error {
		// 	if kh == "" {
		// 		kh = filepath.Join(homeDir, ".ssh", "known_hosts")
		// 	}
		// 	hostKeyCallback, err := knownhosts.New(kh)
		// 	if err != nil {
		// 		return fmt.Errorf("could not create hostkeycallback function (%s) : %v", kh, err)
		// 	}
		// 	hErr := hostKeyCallback(host, remote, pubKey)
		// 	if errors.As(hErr, &keyErr) && len(keyErr.Want) > 0 {
		// 		log.Println("WARNING: %v is not a key of %s, either a MiTM attack or %s has reconfigured the host pub key.", string(pubKey.Marshal()), host, host)
		// 		return keyErr
		// 	} else if errors.As(hErr, &keyErr) && len(keyErr.Want) == 0 {
		// 		log.Println("WARNING: %s is not trusted, adding this key: %q to known_hosts file.", host, string(pubKey.Marshal()))
		// 		return addHostKey(host, remote, pubKey)
		// 	}
		// 	log.Println("Pub key exists for %s.", host)
		// 	return nil
		// }),
	}

	// finally dialing the connection
	client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", input.Host, input.Port), config)
	return
}

func outputListener(stdout, stderr io.Reader, outChan, errChan chan []byte, out chan<- []byte) {
	go func(outChan chan<- []byte) {
		defer close(outChan)
		log.Println("StdOut started")
		scanner := bufio.NewScanner(stdout)
		for {
			if tkn := scanner.Scan(); tkn {
				rcv := scanner.Bytes()
				log.Println("StdOut:", rcv)
				outChan <- rcv
			} else {
				if scanner.Err() != nil {
					outChan <- []byte(scanner.Err().Error())
					log.Println("StdOut:", scanner.Err().Error())
				} else {
					log.Println("StdOut: io.EOF")
				}
				break
			}
		}
		log.Println("StdOut exited")
	}(outChan)

	go func(errChan chan<- []byte) {
		defer close(errChan)
		log.Println("StdErr started")
		scanner := bufio.NewScanner(stderr)
		for {
			if tkn := scanner.Scan(); tkn {
				rcv := scanner.Bytes()
				// raw := make([]byte, len(rcv))
				// copy(raw, rcv)
				log.Println("StdErr:", rcv)
				errChan <- rcv
			} else {
				if scanner.Err() != nil {
					errChan <- []byte(scanner.Err().Error())
					log.Println("StdErr:", scanner.Err().Error())
				} else {
					errChan <- []byte(scanner.Text())
					log.Println("StdErr: io.EOF")
				}
				break
			}
		}
		log.Println("StdErr exited")
	}(errChan)

	go func(outChan, errChan <-chan []byte, out chan<- []byte) {
		defer close(out)
		log.Println("Listener input ...")
		var outOk, errOk bool
		var out1, err1 []byte
		for {
			select {
			case out1, outOk = <-outChan:
				if outOk {
					out <- out1
				}
			case err1, errOk = <-errChan:
				if errOk {
					out <- err1
				}
			}
			if (!outOk) && (!errOk) {
				break
			}
		}
		log.Println("Listener output ...")
	}(outChan, errChan, out)
}

// func checkKnownHosts() ssh.HostKeyCallback {
// 	createKnownHosts()
// 	kh, e := knownhosts.New(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
// 	return kh
// }

// func createKnownHosts() {
// 	f, fErr := os.OpenFile(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"), os.O_CREATE, 0600)
// 	f.Close()
// }

// func addHostKey(host string, remote net.Addr, pubKey ssh.PublicKey) error {
// 	// add host key if host is not found in known_hosts, error object is return, if nil then connection proceeds,
// 	// if not nil then connection stops.
// 	khFilePath := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")

// 	f, fErr := os.OpenFile(khFilePath, os.O_APPEND|os.O_WRONLY, 0600)
// 	if fErr != nil {
// 		return fErr
// 	}
// 	defer f.Close()

// 	knownHosts := knownhosts.Normalize(remote.String())
// 	_, fileErr := f.WriteString(knownhosts.Line([]string{knownHosts}, pubKey))
// 	return fileErr
// }
