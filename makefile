BUILD_DIR=out

all: build run

build: 
	go build -o $(BUILD_DIR)/main main.go

run: 
	./$(BUILD_DIR)/main