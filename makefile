.PHONY: gpush, testing, vet, lint

gpush: # make gpush dir=(directory to push) msg=(message to commit)
	@test -n "$(msg)" || (echo "Usage: make gpush dir=... msg=\"...\"" && exit 1)
	cd backend/ && \
	git add go.mod go.sum $(dir) && \
	git commit -m "$(msg)" && \
	git push

testing:
	cd backend/ && \
	go fmt ./... && \
	

vet:
	cd backend && go vet ./...

lint:
	cd backend && golangci-lint run