submit:
    codecrafters submit
     
run *args: build
    ./bin/app {{args}}

build:
    go build -o ./bin/app ./app

test:
    go test ./...
