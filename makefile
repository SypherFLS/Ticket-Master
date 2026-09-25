.PHONY: gpush, testing, bup, bups, lint, down, callv

gpush: # make gpush dir=(directory to push) msg=(message to commit)
	@test -n "$(msg)" || (echo "Usage: make gpush dir=... msg=\"...\"" && exit 1)
	cd backend/ && \
	git add go.mod go.sum $(dir) && \
	git commit -m "$(msg)" && \
	git push

bups: # пересобрать бекенд
	cd backend && docker build -t tm . && docker-compose up -d --build 'backend'

bup: # собрать весь проект
	cd backend && docker build -t tm . && docker-compose up -d --build 

down: 
	docker-compose down -v

callv:
	cd backend && go-callvis ./cmd/main.go

lint:
	cd backend && golangci-lint run