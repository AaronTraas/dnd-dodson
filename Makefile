run:
	go run main.go restcontrol.go

build:
	go build -o hexmap *.go

dep:
	go get sigs.k8s.io/yaml
