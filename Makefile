BIN_DIR=bin
BINARY_NAME=pokedex-cli

build:
	go build -o ${BIN_DIR}/${BINARY_NAME} main.go

run: build
	./${BIN_DIR}/${BINARY_NAME}

clean:
	go clean
	rm ${BIN_DIR}/${BINARY_NAME}
