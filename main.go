package main

import (
	"context"
	"fmt"
	"log"
)

func main() {
	ctx := context.Background()

	// Get the current external IP
	fmt.Println("Extracting external IP address")

	getter := IPGetter{}
	ip, err := getter.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("IP discovered as " + ip)

	// Update DNS to point at it
	fmt.Println("Comparing IP against DigitalOcean")
	updater := NewDigitalOceanUpdater()

	remoteIP, err := updater.GetCurrent(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DigitalOcean IP is " + remoteIP)

	if remoteIP == ip {
		fmt.Println("IP Addresses match. Skipping update of DNS record")
		return
	}

	err = updater.Update(ctx, ip)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Done! New IP address for %s.%s set to %s\n", subdomain, domain, ip)
}
