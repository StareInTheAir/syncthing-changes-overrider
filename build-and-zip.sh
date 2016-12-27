#!/bin/bash
PROJECT_PATH="github.com/StareInTheAir/syncthing-changes-overrider"
BINARY_NAME="syncthing-changes-overrider"
BINARY_FOLDER="$GOPATH/src/$PROJECT_PATH/bin"

function print_exec_name() {
  printf "$BINARY_FOLDER/${GOOS}_$GOARCH/$BINARY_NAME"
  if [ "$GOOS" == "windows" ]; then
   printf ".exe"
  fi
}

function build() {
  printf "${GOOS}_$GOARCH: building"
  # -ldflags -w = no debug info
  go build -o $(print_exec_name) -ldflags -w "$PROJECT_PATH"
}

function package() {
  printf ", zipping\n"
  zip -qj9 "$BINARY_FOLDER/$BINARY_NAME-${GOOS}_$GOARCH.zip" $(print_exec_name) "$GOPATH/src/$PROJECT_PATH/OverriderConfig-default.json"
}

if [ -d "$BINARY_FOLDER" ]; then
  rm -r "$BINARY_FOLDER"
fi

# Mac
export GOOS=darwin  GOARCH=386;    build && package
export GOOS=darwin  GOARCH=amd64;  build && package

# Linux
export GOOS=linux   GOARCH=386;    build && package
export GOOS=linux   GOARCH=amd64;  build && package
export GOOS=linux   GOARCH=arm;    build && package
export GOOS=linux   GOARCH=arm64;  build && package

# Windows
export GOOS=windows GOARCH=386;    build && package
export GOOS=windows GOARCH=amd64;  build && package
