package config

import (
	"fmt"
	"path"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	configPath := "test_config.yml"
	path, error := filepath.Abs(path.Join("assets", configPath))

	if error != nil {
		t.Fatal("Error: ", error)
	}

	config, error := ParseConfig(path)

	if error != nil {
		t.Fatal("Error: ", error)
	}

	fmt.Println("Raw Config ----------------")
	fmt.Printf("%+v\n", config)

	//j, _ := json.MarshalIndent(config.Plans(), "", "  ")
	fmt.Println("Raw Plans ----------------")
	fmt.Printf("%+v\n", config.Plans())
	//
	//j, _ = json.MarshalIndent(config.Services(), "", "  ")
	fmt.Println("Raw Service ----------------")
	fmt.Printf("%+v\n", config.Services())
}
