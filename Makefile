keygen:
	ssh-keygen -t rsa -b 4096

b:
	go build -o sshchecker main.go

r:
	go run main.go

run:
	./sshchecker -f samples/input.json -o samples/output.json

release:
	goreleaser release --snapshot --clean 

build:
	goreleaser build --single-target --snapshot --clean

