package main

import (
	"github.com/snubwoody/orchestrator/compute"
	"github.com/snubwoody/orchestrator/zones"
)

func main() {

	compute.DeleteInstance(zones.UsEast5a, "my-instance")
}
