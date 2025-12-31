APP_NAME := gosible
BUILD_DIR := build
BINARY := $(BUILD_DIR)/$(APP_NAME)
INSTALL_PATH := /usr/local/bin/$(APP_NAME)

.PHONY: all build clean install uninstall

all: build

build:
	@echo "Compiling $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY) ./cmd/main.go
	@echo "Binary created: $(BINARY)"

install: build
	@echo "Installing in $(INSTALL_PATH)"

	@if [ -f "$(INSTALL_PATH)" ]; then \
		printf "File already exists. Overwrite ? [y/N] "; \
		read confirm; \
		if [ "$$confirm" != "y" ] && [ "$$confirm" != "Y" ]; then \
			echo "aborted"; \
			exit 1; \
		fi; \
	fi

	@if [ ! -w "$(INSTALL_PATH)" ] && [ "$$(id -u)" -ne 0 ]; then \
		sudo cp $(BINARY) $(INSTALL_PATH); \
		sudo chmod +x $(INSTALL_PATH); \
	else \
		cp $(BINARY) $(INSTALL_PATH); \
		chmod +x $(INSTALL_PATH); \
	fi

	@echo "done"

uninstall:
	@echo "Removing $(INSTALL_PATH)..."
	@if [ -f "$(INSTALL_PATH)" ]; then \
		if [ "$$(id -u)" -ne 0 ]; then \
			sudo rm "$(INSTALL_PATH)"; \
		else \
			rm "$(INSTALL_PATH)"; \
		fi; \
		echo "removed"; \
	else \
		echo "No binary found."; \
	fi

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "done"
