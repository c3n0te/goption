.PHONY: option
option:
	go build -ldflags "-s -w" -o ./bin ./option
