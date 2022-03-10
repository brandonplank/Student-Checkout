all: build

build:
	@echo "Building server for GNU/Linux"
	@echo "Building Web Assembly file"
	@GOOS=js GOARCH=wasm go build -o Public/assets/wasm/asm.wasm bin/webasm/main.go 
	@echo "Building binary for the server"
	@go build -v
	@echo "Making the binary executable"
	@chmod +x checkout
run:
	@GOOS=js GOARCH=wasm go build -o Public/assets/wasm/asm.wasm bin/webasm/main.go
	@go run .