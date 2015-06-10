#!/bin/bash

set -e

echo "mode: set" > acc.out
for Dir in $(find ./* -maxdepth 10 -type d );
do
        if ls $Dir/*.go &> /dev/null;
        then
            go test -v -bench=. -coverprofile=profile.out $Dir
            if [ -f profile.out ]
            then
                cat profile.out | grep -v "mode: set" >> acc.out
            fi
fi
done
gocov convert acc.out > acc.json
gocov annotate acc.json
goveralls -gocovdata=acc.json
rm -rf ./profile.out
rm -rf ./acc.out
rm -rf ./acc.json
