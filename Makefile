.PHONY: all clean deps build volume wifi iface openconnect battery install dir copy

all: clean build

clean:
	rm -f volume wifi iface openconnect battery

deps:
	go get -t ./...

build: volume wifi iface openconnect battery

volume:
	CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' ./cmd/volume/

wifi:
	CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' ./cmd/wifi/

iface:
	CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' ./cmd/iface/

openconnect:
	CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' ./cmd/openconnect/

battery:
	CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' ./cmd/battery/

install: build dir copy

dir:
	mkdir -p ~/.config/i3blocks/blocklets/go

copy:
	cp volume ~/.config/i3blocks/blocklets/go
	cp wifi ~/.config/i3blocks/blocklets/go
	cp iface ~/.config/i3blocks/blocklets/go
	cp openconnect ~/.config/i3blocks/blocklets/go
	cp battery ~/.config/i3blocks/blocklets/go
