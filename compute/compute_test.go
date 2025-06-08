package compute

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/snubwoody/orchestrator/zones"
)

func TestInsertInstance(t *testing.T) {
	zone := zones.UsEast5a
	name := fmt.Sprintf("test-instance-%s", uuid.New())
	client, err := NewClient()

	if err != nil {
		t.Errorf("Error creating compute client: %s", err)
	}

	client.InsertInstance(name, zone)

	instances, err := client.ListInstances(zone)

	_, err = client.DeleteInstanceAsync(name, zone)

	if err != nil {
		t.Errorf("Error deleting instance: %s", err)
	}

	// FIXME search for the name instead of using len
	if len(instances) != 1 {
		t.Errorf("Instance not created, instances lenght: %v", len(instances))
	}
}

func TestDeleteInstance(t *testing.T) {
	zone := zones.UsEast5a
	name := fmt.Sprintf("test-instance-%s", uuid.New())
	client, err := NewClient()

	if err != nil {
		t.Errorf("Error creating compute client: %s", err)
	}

	client.InsertInstance(name, zone)
	client.DeleteInstance(name, zone)

	instances, err := client.ListInstances(zone)
	if err != nil {
		t.Errorf("Error listing instances: %s", err)
	}

	if len(instances) != 0 {
		t.Error("Instance not cleaned up")
	}
}
