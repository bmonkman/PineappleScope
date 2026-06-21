run:
	go run cmd/pineapplescope/main.go

build: 
	go build ./cmd/pineapplescope

# Self-contained multi-stage build: compiles the CGO binary inside Docker,
# so no local Go/xgo toolchain or cross-compilation is needed.
container:
	docker build --build-arg VERSION=$(shell git rev-parse --short HEAD) -t bmonkman/pineapplescope -f Dockerfile .

push:
	docker push bmonkman/pineapplescope