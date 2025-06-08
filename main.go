package main

import (
	"github.com/snubwoody/orchestrator/compute"
	"github.com/snubwoody/orchestrator/zones"
)

func main() {
	compute.ListInstances(zones.UsEast5a)
}
