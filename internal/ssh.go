package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
)

func doSSH(ind int, input Input, rsa, kh string) (results []string, errors []string) {

	// dialing the ssh connection
	client, err := dialConn(input, rsa, kh)
	if err != nil {
		log.Println(ind, input.Host, "Connection Error", err.Error())
		errors = append(errors, err.Error())
		return
	}
	defer client.Close()

	// sshCompleted := make(chan struct{})
	// total_commands := len(input.Commands)
	for _, cmd := range input.Commands {
		log.Println(ind, input.Host, "Running cmd ---", cmd)

		// creating command based session and closing old once completed
		session, err := client.NewSession()
		if err != nil {
			log.Println(ind, input.Host, "Session Error", cmd, err.Error())
			errors = append(errors, err.Error())
			break
		}

		// executing the command
		out, err := session.CombinedOutput(cmd)
		if err != nil {
			log.Println(ind, input.Host, "Command Error", cmd, err.Error())
			errors = append(errors, err.Error())
			session.Close()
			break
		}
		log.Println(ind, input.Host, "Running cmd completed ---", cmd)
		results = append(results, string(out))
		session.Close()

	}
	return
}

func dialConn(input Input, f, kh string) (client *ssh.Client, err error) {
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

// func putDataToChan(ch chan<- []byte, read io.Reader, t string) {
// 	defer close(ch)
// 	log.Println(t, "started")
// 	scanner := bufio.NewScanner(read)
// 	for scanner.Scan() {
// 		rcv := scanner.Bytes()
// 		log.Println(t, ":", string(rcv))
// 		ch <- rcv
// 	}
// 	if err := scanner.Err(); err != nil {
// 		ch <- []byte(scanner.Err().Error())
// 		log.Println(t, ":", scanner.Err().Error())
// 	} else {
// 		log.Println(t, ": io.EOF")
// 	}
// 	log.Println(t, " exited")
// }

// func outputListener(s string, outCh <-chan []byte, errCh <-chan []byte) (results []string, errors []string) {
// 	var o, e []byte
// 	outOk, errOk := true, true
// 	log.Println(s, "Listener entering ...")
// 	for {
// 		select {
// 		case o, outOk = <-outCh:
// 			if outOk {
// 				results = append(results, string(o))
// 			}
// 		case e, errOk = <-errCh:
// 			if errOk {
// 				errors = append(errors, string(e))
// 			}
// 		}
// 		if (!outOk) && (!errOk) {
// 			break
// 		}
// 	}
// 	log.Println(s, "Listener exit")
// 	return
// }

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
