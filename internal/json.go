package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type _json struct {
	ifile string
	ofile string
}

func NewJSON(ifile, ofile string) SSHChecker {
	return &_json{ifile, ofile}
}

func (c *_json) GetInput() (out []Input) {
	fileData, err := os.ReadFile(c.ifile)
	if err != nil {
		log.Fatal("Error while read the json file", err)
	}

	err = json.Unmarshal(fileData, &out)
	if err != nil {
		log.Fatal("Error while converting to json", err)
	}

	return out
}

func (c *_json) ProduceOutput(ch <-chan Output, exitCh chan<- struct{}) {
	log.Println("initiate the read output...")
	var out []Output
	var counter int = 0
	for r := range ch {
		out = append(out, r)
		counter++
		log.Println(counter, "reading output...", r.I.Host)
	}

	log.Println("data received", len(out), "| going to write as file -", c.ofile)
	outJson, err := json.Marshal(out)
	if err != nil {
		log.Fatal("Error while marshaling to json", err)

	} else {
		err = os.WriteFile(fmt.Sprintf("%s", c.ofile), outJson, 0644)
		if err != nil {
			log.Fatal("Error while writing into file", err)
		}
	}
	exitCh <- struct{}{}
}
