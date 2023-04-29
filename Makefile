deps:
	go mod download
	go mod verify
	go mod tidy

start:
	go run main.go --config config.yaml --dry

check:
	go run main.go --config config.yaml

learn:
	go run main.go --config config.yaml --learn

build: deps
	go build -o build/stukach main.go
	upx build/stukach

install:
	cp build/stukach /usr/bin/stukach
	chown root:root /usr/bin/stukach
	chmod 0755 /usr/bin/stukach
	restorecon /usr/bin/stukach

uninstall:
	rm /usr/bin/stukach
