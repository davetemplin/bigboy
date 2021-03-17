#!/bin/bash
trap "exit" ERR

OS=$1

build_darwin () {
  GOOS=darwin 
  GOARCH=amd64
  name=bigboy-${GOOS}-${GOARCH}
  GOOS=$GOOS GOARCH=$GOARCH go build -o bin/$name
  tar -czvf bin/$name.tar.gz bin/$name
}

build_linux () {
  GOOS=linux 
  GOARCH=amd64
  name=bigboy-${GOOS}-${GOARCH}
  GOOS=$GOOS GOARCH=$GOARCH go build -o bin/$name
  tar -czvf bin/$name.tar.gz bin/$name
}

build_windows() {
  GOOS=windows 
  GOARCH=amd64
  name=bigboy-${GOOS}-${GOARCH}
  GOOS=$GOOS GOARCH=$GOARCH go build -o bin/$name.exe
  zip bin/$name bin/$name.exe
}

if [[ -z $OS ]]; then
  build_darwin
  build_linux
  build_windows
else 
  case $OS in
    darwin)
      build_darwin
      ;;
    linux)
      build_linux
      ;;
    windows)
      build_windows
      ;;
    *)
      echo 'Invalid OS name. Has to be one of darwin|linux|windows'
      exit 1
      ;;
  esac
fi

