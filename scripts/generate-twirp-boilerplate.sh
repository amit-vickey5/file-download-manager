cd $GOPATH"/src/github.com/amit/file-download-manager"
protoc --proto_path=$GOPATH/src:. --twirp_out=. --go_out=. ./rpc/service.proto