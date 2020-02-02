default:
	go build

arm:
	GOOS=linux GOARCH=arm GOARM=5 go build

deploy: arm
	chmod +x ./home-updater

	rsync ./home-updater ./.env pi@10.0.1.2:/tmp/
	ssh pi@10.0.1.2 sudo mv /tmp/home-updater /usr/local/bin/home-updater
	ssh pi@10.0.1.2 sudo mv /tmp/.env /usr/local/bin/home-updater.env
	ssh pi@10.0.1.2 sudo chown root:root /usr/local/bin/home-updater*

	nomad job run -address="http://10.0.1.2:4646" nomad.hcl