#!/bin/bash

[ -d ./bin ] || {
    echo "./bin folder not existent"
    exit 1
}

cp -r ./web ./bin/web
cd ./bin

for f in ./*; do
    if [ -f $f ]; then
        sha256sum $f >> ./sums.txt
    fi
done

tar -czvf ./release.tar.gz *

sha256sum ./release.tar.gz > ./sum.txt

for f in ./*; do
    if [ "$f" != "./release.tar.gz" ] && [ "$f" != "./sum.txt" ]; then
        rm -r -f $f
    fi
done