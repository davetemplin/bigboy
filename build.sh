#!/bin/bash
trap "exit" ERR

OS=$1

build_mac () {
  GOOS=darwin 
  GOARCH=amd64
  go build -o bin/mac/bigboy
}

build_linux () {
  GOOS=linux 
  GOARCH=amd64
  go build -o bin/linux/bigboy
}

build_windows() {
  GOOS=windows 
  GOARCH=amd64
  go build -o bin/windows/bigboy.exe
}

if [[ -z $OS ]]; then
  build_mac
  build_linux
  build_windows
else 
  case $OS in
    mac)
      build_mac
      ;;
    linux)
      build_linux
      ;;
    windows)
      build_windows
      ;;
    *)
      echo 'Invalid OS name. Has to be one of mac|linux|windows'
      exit 1
      ;;
  esac
fi

