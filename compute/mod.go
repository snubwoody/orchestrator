package compute

import (
	"context"
	"fmt"
	"log"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/snubwoody/orchestrator/config"
	"google.golang.org/api/iterator"
)

// List all the active instances in a zone
func ListInstances(zone string) {
	ctx := context.Background()
	fmt.Println("Creating client")

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		fmt.Printf("NewInstancesRESTClient error: %s", err)
	}

	defer instancesClient.Close()

	req := &computepb.ListInstancesRequest{
		Project: config.ProjectId,
		Zone:    zone,
	}

	fmt.Println("Fetching instances")
	it := instancesClient.List(ctx, req)
	count := it.PageInfo().Remaining()

	if count == 0 {
		fmt.Println("No instances found in zone")
		return
	}

	for {
		instances, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			fmt.Printf("Iter err: %s", err)
			break
		}

		fmt.Printf("%s %s\n", instances.GetName(), instances.GetMachineType())
	}
}

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
