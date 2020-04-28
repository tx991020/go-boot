all: deps


deps:
#	GOOS=windows go build -o WebGenerator.exe main.go
#	GOOS=linux go build -o  WebGenerator main.go
	GOOS=darwin go build -o  WebGenerator_mac main.go