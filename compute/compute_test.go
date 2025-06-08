package compute

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/snubwoody/orchestrator/zones"
)

func TestDeleteInstance(t *testing.T) {
	zone := zones.UsEast5a
	name := fmt.Sprintf("test-instance-%s", uuid.New())

	InsertInstance(name, zone)
	DeleteInstance(name, zone)

	instances, err := ListInstances(zone)

	if err != nil {
		t.Errorf("Error listing instances: %s", err)
	}

	if len(instances) != 0 {
		t.Error("Instance not cleaned up")
	}
}
