BUILD_DIR=$(PWD)/bin
SCHEMA_DIR=../../infrastructure/schema/
build:
	echo "Building profile service"
	mkdir -p ${BUILD_DIR} && GOOS=linux GOARCH=386 go build -o ${BUILD_DIR}/profile-service main.go

test-unit:
	echo "Running unit tests"
	go test -short ./...

test-integration:
	echo "Running integration tests"
	docker build . -t profile-service \
		&& mkdir -p test_data \
		&& grep -v '^--*' ${SCHEMA_DIR}/schema.sql >test_data/schema.sql \
		&& go test --tags=integration ./...

test-all: test-unit test-integration

clean:
	echo "Running cleanup"
	rm bin/*
	docker rmi -f profile-service
