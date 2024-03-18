build:
	go build -o validator-duties main.go

clean: 
	go clean
	rm validator-duties
