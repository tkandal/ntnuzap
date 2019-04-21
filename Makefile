
ntnuzap: $(wildcard *.go)
	dep ensure -v
	go build -v .

vet:
	go vet .

lint:
	golint .

graph:
	dep status -dot | dot -T png | open -f -a /Applications/Preview.app

graphwin:
	dep status -dot | dot -T png -o status.png; start status.png

graphlin:
	dep status -dot | dot -T png -o status.png ; xdg-open status.png

clean:
	rm -rf vendor
	rm -f status.png

.PHONY: ntnuzap lint vet

