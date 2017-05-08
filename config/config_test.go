package config

import (
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	configPath := "test_config.yml"
	path, err := filepath.Abs(path.Join("assets", configPath))

	if err != nil {
		t.Fatal("Error: ", err)
	}

	config, err := ParseConfig(path)

	if err != nil {
		t.Fatal("Error: ", err)
	}

	fmt.Println("Raw Config ----------------")
	fmt.Printf("%+v\n", config)

	//j, _ := json.MarshalIndent(config.Plans(), "", "  ")
	fmt.Println("Raw Plans ----------------")
	fmt.Printf("%+v\n", config.Plans())
	//
	//j, _ = json.MarshalIndent(config.Services(), "", "  ")
	//fmt.Println("Raw Service ----------------")
	//fmt.Printf("%+v\n", config.Services())

	j, _ := json.MarshalIndent(config.Services(), "", "  ")
	fmt.Println("Raw Service ----------------")
	fmt.Printf("%+s\n", j)
}
