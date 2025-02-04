package main

import (
	"fmt"
	"log"
	"os"

	"cf-application-discovery/pkg/discover"

	"gopkg.in/yaml.v3"
)

func main() {

	// Check if a file path is provided as an argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path_to_manifest.yml>")
		return
	}

	var manifestFilePath = os.Args[1]

	// Read the YAML file
	data, err := os.ReadFile(manifestFilePath)
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		return
	}

	// Unmarshal the YAML data into the Manifest struct
	var cfApplications discover.Manifest
	err = yaml.Unmarshal(data, &cfApplications)
	if err != nil {
		fmt.Printf("Error unmarshalling YAML: %v\n", err)
		return
	}

	if len(cfApplications.Applications) > 0 {
		for i, v := range cfApplications.Applications {
			d, err := discover.Discover(*v, cfApplications.Version, "")
			if err != nil {
				log.Fatal(err)
			}
			m, err := yaml.Marshal(d)
			if err != nil {
				fmt.Printf("Error marshalling YAML file: %v\n", err)
				return
			}
			fmt.Printf("#%d\n%s\n", i, m)
		}
	} else {
		fmt.Println("No applications found.")
	}

}
