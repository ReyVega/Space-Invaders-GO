build:
	go build main.go

test: build
	./main
	
clean:
	rm -rf main