build:
	go build -o validator-duties cmd/validator-duties/main.go

clean: 
	go clean
	rm validator-duties
