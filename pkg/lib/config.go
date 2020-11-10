package lib

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	CommitsCount   int      `json:"commitsCount"`
	DiffFilter     string   `json:"diffFilter"`
	DirectoryNames []string `json:"directoryNames"`
	GitHashNewest  string   `json:"gitHashNewest"`
	GitHashOldest  string   `json:"gitHashOldest"`
	OutputPath     string   `json:"outputPath"`
	RepositoryPath string   `json:"repositoryPath"`
}

func LoadConfiguration(file string) Configuration {
	var config Configuration
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		fmt.Println("[ERROR] Could not decode configuration file")
		fmt.Println(err.Error())
	}

	return config
}
