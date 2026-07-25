BINARY_NAME=raw_api

start_db:
	podman-compose up -d

stop_db:
	podman-compose down

build: start_db
	@echo "Building..."
	@go build -v -o ${BINARY_NAME} ./cmd/server
	@echo "Built!"

run: start_db build
	@echo "Starting..."
	@env ./${BINARY_NAME} &
	@echo "Started!"

clean: 
	@echo "Cleaning..."
	@go clean
	@rm ${BINARY_NAME}
	@echo "Cleaned!"

start: run

stop: stop_db
	@echo "Stopping..."
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
	@echo "Stopped!"

restart: stop stop_db start_db start

test:
	gotest -count=1 ./... -v
