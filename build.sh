#!/bin/bash

echo "prepare build target dir"
### prepare build target dir
if [ ! -d "./build" ]; then
  mkdir build
fi

if [ ! -d "./build/bin" ]; then
  mkdir ./build/bin
fi


if [ ! -d "./build/conf" ]; then
  mkdir ./build/conf
fi

if [ ! -d "./build/log" ]; then
  mkdir ./build/log
fi

### version
echo "set version"

package_name="mediachop"

Path="${package_name}/config"

echo "-version    $Path"

GitCommit=$(git rev-parse --short HEAD || echo unsupported)
GoVersion=$(go version)
BuildTime=$(date "+%Y-%m-%d %H:%M:%S")

echo "GoVersion    $GoVersion"
echo "GitCommit    $GitCommit"
echo "BuildTime    $BuildTime"

## build
echo "build ..."
go build -ldflags "-X '$Path.GoVersion=$GoVersion' -X '$Path.GitCommit=$GitCommit' -X '$Path.BuildTime=$BuildTime'"

## mv build file
cp mediachop ./build/bin
cp ./tools/*.sh ./build/bin
cp ./conf/mediachop.yaml ./build/conf

echo "package ..."
## package
cd build
tar -czvf mediachop.tar.gz --exclude=log/*.log --exclude=log/*.gz ./
mv mediachop.tar.gz ../

echo "done ..."