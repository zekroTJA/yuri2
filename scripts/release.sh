#!/bin/bash

OSES=(
    "linux" 
    "windows" 
    "darwin" 
    "freebsd"
)

ARCHES=(
    "386 amd64 arm"
    "386 amd64"
    "386 amd64"
    "386 amd64"
)



# -------------------------------------------------

function build {
    (env GOOS=$1 GOARCH=$2 \
        make BINPATH="./_release" build)
}

# -------------------------------------------------

[ -f ./Makefile ] || {
    echo "[ FATAL ] Makefile was not found"
    exit
}

[ -d "./_release/web" ] || {
    mkdir -p ./_release/web
}

[ -d "./release" ] || {
    mkdir ./release
}

cp -r ./web/* ./_release/web

make deps

i=0
for OS in ${OSES[*]}; do
    _ARCHES=${ARCHES[i]}
    for ARCH in ${_ARCHES[*]}; do
        build $OS $ARCH
    done
    i=$((i + 1))
done

for f in ./_release/*; do
    [ -f $f ] && {
        sha256sum $f >> ./_release/sums.txt
    }
done

cd ./_release
tar -czvf ../release/release.tar.gz * \
    --exclude=*.gz
cd ..

sha256sum ./release/release.tar.gz > ./release/sum.txt

rm -r -f ./_release
