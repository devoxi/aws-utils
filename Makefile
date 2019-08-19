all: build

build:
	go build -o bin/get-ec2-ip get_ec2_ip.go
	go build -o bin/get-emr-master-ip get_emr_master_ip.go
