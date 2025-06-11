package test

import (
	"github.com/snubwoody/orchestrator/compute"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestParseYaml(t *testing.T) {
	data := map[string]any{
		"type":        "PERSISTENT",
		"mode":        "READ_WRITE",
		"bool":        true,
		"auto-delete": false,
		"disk-size":   5,
	}

	yamlData, err := yaml.Marshal(data)

	if err != nil {
		t.Errorf("yaml marshal fail: %s", err)
	}

	disk := compute.Disk{}
	err = yaml.Unmarshal(yamlData, disk)
	if err != nil {
		t.Errorf("Failed to parse disk yaml: %s", err)
	}

}
