dist:
	mkdir -p dist
	minify web/index.js web/*.html web/*.css -o dist/
	cp -rf web/*.png dist/web/
	cp -f web/axios.min.js dist/web/axios.min.js
	cp -f web/vue.global.prod.js dist/web/vue.js
	cp -f web/vue-router.global.prod.js dist/web/vue-router.js

linux:
	go build

windows:
	CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build
