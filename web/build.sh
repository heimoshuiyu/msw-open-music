rm -rf build
cp -raf public build
./node_modules/.bin/esbuild src/index.jsx --bundle --outfile=build/msw-open-music.js --alias:react=preact/compat --alias:react-dom=preact/compat --minify --analyze
cat public/index.html | sed "s/%PUBLIC_URL%/$PUBLIC_URL/" > build/index.html

echo "Build done, output files under ./build directory"
