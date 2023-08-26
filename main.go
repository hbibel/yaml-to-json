package main

import (
	"bufio"
	"fmt"
	"hbibel/yaml-to-json/json"
	"hbibel/yaml-to-json/yaml"
	"log"
	"os"
)

type Config struct {
}

func main() {
	yamlFilePath := "test.yaml"
	jsonFilePath := "test.json"
	var err error

	_, err = os.Stat(jsonFilePath)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	} else if err == nil {
		err = os.Remove(jsonFilePath)
		if err != nil {
			log.Fatal(err)
		}
	}

	yamlFile, err := os.Open(yamlFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer yamlFile.Close()

	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	var tokens chan yaml.Token = make(chan yaml.Token)
	var lines chan string = make(chan string)
	yaml.Tokenize(lines, tokens)
	events := yaml.TokensToEvents(tokens)
	jsonChunks := json.RenderEvents(events)

	outDone := make(chan bool)
	go func() {
		writer := bufio.NewWriter(jsonFile)
		for chunk := range jsonChunks {
			fmt.Fprint(writer, chunk)
		}
		writer.Flush()
		outDone <- true
	}()

	scanner := bufio.NewScanner(yamlFile)
	for scanner.Scan() {
		lines <- scanner.Text()
	}
	close(lines)
	<-outDone
}
