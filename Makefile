dist:
	mkdir -p dist
	minify web/index.js web/*.html web/*.css -o dist/
	cp -rf web/*.png dist/web/
	cp -f web/axios.min.js dist/web/axios.min.js
	cp -f web/vue.global.prod.js dist/web/vue.js
	cp -f web/vue-router.global.prod.js dist/web/vue-router.js
