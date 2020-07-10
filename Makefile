all: build

build:
	wire /src/cmd/panel
	go build -o panel /src/cmd/panel/
