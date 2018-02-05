.PHONY: all clean deps build volume wifi iface openconnect battery install dir copy

all: clean build

clean:
	rm -f volume wifi iface openconnect battery

deps:
	go get -t ./...

build: volume wifi iface openconnect battery

volume:
	CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o dist/volume ./cmd/volume/

wifi:
	CGO_ENABLED=0 go build -ldflags '-extldflags "-static"'  -o dist/wifi ./cmd/wifi/

iface:
	CGO_ENABLED=0 go build -ldflags '-extldflags "-static"'  -o dist/iface ./cmd/iface/

openconnect:
	CGO_ENABLED=0 go build -ldflags '-extldflags "-static"'  -o dist/openconnect ./cmd/openconnect/

battery:
	CGO_ENABLED=0 go build -ldflags '-extldflags "-static"'  -o dist/battery ./cmd/battery/

install: build dir copy

dir:
	mkdir -p ~/.config/i3blocks/blocklets/go

copy:
	cp dist/volume ~/.config/i3blocks/blocklets/go
	cp dist/wifi ~/.config/i3blocks/blocklets/go
	cp dist/iface ~/.config/i3blocks/blocklets/go
	cp dist/openconnect ~/.config/i3blocks/blocklets/go
	cp dist/battery ~/.config/i3blocks/blocklets/go
