INPUT_DIR := cmd/server
OUT_DIR := bin
OUT_BIN := $(OUT_DIR)/main

build: $(INPUT_DIR)/main.go
	go build -o ./$(OUT_BIN) ./$(INPUT_DIR)/main.go

run: build
	./$(OUT_BIN)

clean:
	rm -rf ./$(OUT_DIR)/*