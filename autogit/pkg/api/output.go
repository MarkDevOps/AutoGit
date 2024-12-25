package api

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func WriteOutput(output interface{}, filePath string) error {
	yamlOutput, err := yaml.Marshal(output)
	if err != nil {
		return fmt.Errorf("error marshalling output: %w", err)
	}

	if err := os.WriteFile(filePath, yamlOutput, 0644); err != nil {
		return fmt.Errorf("error writing to file %s: %w", filePath, err)
	}

	fmt.Printf("Output written to %s\n", filePath)
	return nil
}
