rm -rf build
cp -raf public build
./node_modules/.bin/esbuild src/index.jsx --bundle --minify --outfile=build/msw-open-music.js --sourcemap
cat public/index.html | sed "s/%PUBLIC_URL%/$PUBLIC_URL/" > build/index.html
