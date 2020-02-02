package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

var domains = []string{
	"home.ryannixon.com",
	"taiidani.com",
	"home.taiidani.com",
	"momentumdance.org",
}

func main() {
	ctx := context.Background()

	fmt.Println("Application Starting")

	for {
		// Get the current external IP
		fmt.Println("Extracting external IP address")

		getter := IPGetter{}
		ip, err := getter.Get(ctx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("IP discovered as " + ip)

		// Update DNS to point at it
		fmt.Println("Comparing IP against DigitalOcean records...")
		updater := NewDigitalOceanUpdater()

		for _, domain := range domains {
			remoteIP, err := updater.GetCurrent(ctx, domain)
			if err != nil {
				log.Fatalf("Error extracting domain for %s: %s", domain, err)
			}
			fmt.Printf("DigitalOcean IP for %s is %s\n", domain, remoteIP)

			if remoteIP == ip {
				fmt.Println("IP Addresses match. Skipping update of DNS record")
			} else {
				err = updater.Update(ctx, domain, ip)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("Done! New IP address for %s set to %s\n", domain, ip)
			}
		}

		fmt.Println("Sleeping for 120 minutes")
		time.Sleep(time.Minute * 120)
	}
}
