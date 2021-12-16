dist:
	cd web && npm install
	cd web && npm run build

linux:
	go build -v

windows:
	CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -v
