deps:
	go mod download
	go mod verify
	go mod tidy

start:
	go run main.go

port_forward:
	ssh -L 127.0.0.1:11333:192.168.47.3:11333 holod.local

