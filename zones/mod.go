package zones

import (
	"context"
	"fmt"
	"log"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/snubwoody/orchestrator/config"
	"google.golang.org/api/iterator"
)

const (
	UsEast5a = "us-east5-a"
)

func ListZones() {
	ctx := context.Background()

	zonesClient, err := compute.NewZonesRESTClient(ctx)

	if err != nil {
		log.Fatalf("Failed to create zones client: %s", err)
	}

	defer zonesClient.Close()

	req := &computepb.ListZonesRequest{
		Project: config.ProjectId,
	}

	it := zonesClient.List(ctx, req)
	fmt.Println("Available zones:")
	for {
		zone, err := it.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Fatalf("Failed to fetch zone: %s", err)
		}

		fmt.Printf("- %s (%s)\n", zone.GetName(), zone.GetStatus())
	}

}
