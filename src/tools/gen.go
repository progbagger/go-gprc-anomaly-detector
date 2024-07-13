package tools

import (
	_ "google.golang.org/grpc"
)

//go:generate protoc --go_out=../generated --go_opt=paths=source_relative --go-grpc_out=../generated --go-grpc_opt=paths=source_relative --gorm_out=engine=postgres:../generated -I $GOPATH/pkg/mod/github.com/infobloxopen/protoc-gen-gorm@v1.1.4/proto/ --gorm_opt=paths=source_relative -I ../proto/ ../proto/types.proto
