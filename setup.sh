#!/bin/bash

command_exists () {
    type "$1" &> /dev/null ;
}

if ! command_exists go; then
    echo "Go is not installed. Please install Go and try again."
    exit 1
fi

if [ ! -f go.mod ]; then
    echo "Initializing Go module..."
    go mod init ProtoZ
fi

# Add required dependencies to go.mod
echo "Adding dependencies..."
go get github.com/chromedp/chromedp
go get github.com/tomnomnom/unfurl

# Compile the Go code
echo "Building ProtoZ..."
go build -o ProtoZ main.go

if [ ! -f ProtoZ ]; then
    echo "Build failed. Exiting."
    exit 1
fi

MODE="search"
JS="window.hacker || window[1337] ? 'Vulnerable' : 'Not Vulnerable'"

while getopts ":f:m:j:" opt; do
  case $opt in
    f) URL_FILE="$OPTARG"
    ;;
    m) MODE="$OPTARG"
    ;;
    j) JS="$OPTARG"
    ;;
    \?) echo "Invalid option -$OPTARG" >&2
        exit 1
    ;;
  esac
done

if [ "$URL_FILE" != "" ]; then
    echo "Running ProtoZ with URL file $URL_FILE and mode $MODE..."
    awk '{print $1}' "$URL_FILE" | ./ProtoZ -j "$JS" -m "$MODE"
else
    echo "Running ProtoZ and expecting URLs from stdin with mode $MODE..."
    awk '{print $1}' | ./ProtoZ -j "$JS" -m "$MODE"
fi

