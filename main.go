package main

import (
	"github.com/snubwoody/orchestrator/compute"
	"github.com/snubwoody/orchestrator/zones"
)

func main() {
	compute.InsertInstance(zones.UsEast5a)
}
