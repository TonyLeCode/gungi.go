.PHONY: run-client run-server

run-server:
	cd ./server && go run ./cmd/server/.

run-client:
	cd ./client && npm run dev