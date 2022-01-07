package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	colorReset := "\033[0m"

	colorRed := "\033[31m"
	colorGreen := "\033[32m"
	colorYellow := "\033[33m"
	colorBlue := "\033[34m"

	files, err := ioutil.ReadDir("./roles")
	if err != nil {
		log.Fatal(err)
	}

	var noDocsServices []string
	var notInIndex []string
	var notUsingIncludes []string

	var serviceCount int
	var happyServices int
	var sadServices int

	fmt.Println("Checking services:")
	fmt.Println()
	for _, file := range files {
		// Filter out HomelabOS internals
		if !strings.Contains(file.Name(), "homelabos") &&
			file.Name() != "tor" &&
			file.Name() != "docs" {
			serviceOk := true
			serviceCount++
			// Detect if the service has a doc file
			if _, err := os.Stat("./docs_software/" + file.Name() + ".md"); errors.Is(err, os.ErrNotExist) {
				noDocsServices = append(noDocsServices, file.Name())
				serviceOk = false
			}

			// Detect if the service is included in docs/index.md
			buffer, err := ioutil.ReadFile("docs/index.md")
			if err != nil {
				panic(err)
			}
			fileContents := string(buffer)

			if !strings.Contains(fileContents, file.Name()) {
				notInIndex = append(notInIndex, file.Name())
				serviceOk = false
			}

			// Detect if the service is using the new include style
			buffer, err = ioutil.ReadFile("roles/" + file.Name() + "/tasks/main.yml")
			if err != nil {
				// File doesn't exist
				serviceOk = false
			}
			fileContents = string(buffer)

			if !strings.Contains(fileContents, "includes/start.yml") {
				notUsingIncludes = append(notUsingIncludes, file.Name())
				serviceOk = false
			}

			// Output service status
			if serviceOk {
				fmt.Print(string(colorGreen), ".")
				happyServices++
			} else {
				fmt.Print(string(colorRed), "X")
				sadServices++
			}
		}
	}

	fmt.Println(string(colorReset))
	fmt.Println()
	fmt.Print("Detected services: ", string(colorBlue))
	fmt.Printf("%d", serviceCount)
	fmt.Println(string(colorReset))
	fmt.Print("Happy services: ", string(colorGreen))
	fmt.Printf("%d", happyServices)
	fmt.Println(string(colorReset))
	fmt.Print("Sad services: ", string(colorRed))
	fmt.Printf("%d", sadServices)
	fmt.Println(string(colorReset))

	fmt.Print("Services without documentation: \n", string(colorYellow))
	for _, serviceName := range noDocsServices {
		fmt.Print(serviceName + "\n")
	}

	fmt.Println(string(colorReset))

	fmt.Print("Services not in documentation index.md: \n", string(colorYellow))
	for _, serviceName := range notInIndex {
		fmt.Print(serviceName + "\n")
	}

	fmt.Println(string(colorReset))

	fmt.Print("Services not using the new includes format:", string(colorYellow))
	for _, serviceName := range notUsingIncludes {
		fmt.Print(serviceName + "\n")
	}
}
