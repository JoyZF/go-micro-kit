#sh
# proto generate shell
# replace api "go-micro.dev/v4/api, client "go-micro.dev/v4/client,server,go-micro.dev/v4/server in *.pb.micro.go before generate success
protoc --proto_path=./proto --micro_out=./proto --go_out=./proto ./proto/*.proto