package usage

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// ApplyYaml Extract usage from Yaml data and applies it to the command hierarchy.
// It prints any errors to stderr.
func ApplyYaml(apply func(Usage), yamlUsage []byte) bool {
	var use Usage
	err := yaml.Unmarshal(yamlUsage, &use)
	if err != nil {
		fmt.Fprintf(os.Stderr, "usage: %v\n", err)
		return false
	}
	apply(use)
	return true
}
