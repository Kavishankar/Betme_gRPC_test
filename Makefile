BIN_DIR = bin
PROTO_DIR = proto
PROTO_FILENAME = betme_test.proto
PROTO_FILE = ${PROTO_DIR}/${PROTO_FILENAME}
PACKAGE = betme_test
SERVER_DIR = server
CLIENT_DIR = client
SERVER_BIN = ${SERVER_DIR}
CLIENT_BIN = ${CLIENT_DIR}

generate:
	protoc --go_opt=M${PROTO_FILE}=${PROTO_DIR}/ --go_out=. --go-grpc_opt=M${PROTO_FILE}=${PROTO_DIR}/ --go-grpc_out=. ${PROTO_DIR}/*.proto

bump: generate
	go get -u ./...

build: 	generate
	go build -o ${BIN_DIR}/${SERVER_BIN} ./${SERVER_DIR}
	go build -o ${BIN_DIR}/${CLIENT_BIN} ./${CLIENT_DIR}

test:
	go test ./...

clean:
	rm -rf ${PROTO_DIR}/*.pb.go
	rm -rf ${BIN_DIR}