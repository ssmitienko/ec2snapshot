ec2snapshot: main.go 
	CGO_ENABLED=0 go build -o ec2snapshot .

docker:
	docker build -t ec2snapshot .

install: ec2snapshot
	mkdir -p /opt/bin
	install -o root -g root -m 0755 ec2snapshot /opt/bin/ec2snapshot
