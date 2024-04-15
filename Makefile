build :
	go build -o bin/salejobs

run: 
	build ./bin/salejobs

test:
	go test -v ./...
