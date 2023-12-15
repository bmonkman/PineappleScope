run:
	go run cmd/pineapplescope/main.go

build: 
	go build ./cmd/pineapplescope

container:
	# Wanted to do this when using the super light scratch container but had to do some other stuff in the container so had to drop it
	#CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
	# Wanted to do this, but SQLite won't cross-compile properly from OSX->Linux with CGO_ENABLED
	#CGO_ENABLED=1 GOOS=linux go build -o main .
	
	go install src.techknowlogick.com/xgo@latest
	xgo --targets=linux/amd64 .
	
	docker build -t bmonkman/pineapplescope -f Dockerfile .

push:
	docker push bmonkman/pineapplescope