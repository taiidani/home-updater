# Home Updater

This simple script reaches out to icanhazip to find your external IP address, then updates the "home.ryannixon.com" DNS record in DigitalOcean with the new value. It is expected to be deployed as a recurring job within the network needing its address updated. Pretty simple DynDNS substitute!

To execute the script:

```
go build
DIGITALOCEAN_TOKEN=1234abcd ./home-updater
```