package main

import (
	"context"
	"fmt"
	"go-web/tls/search"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

type service struct {
	search.UnimplementedSearchServiceServer
}

func (s *service) Search(ctx context.Context, req *search.SearchRequest) (res *search.SearchResponse, err error) {
	fmt.Println("收到客户端：" + req.GetRequest())
	return &search.SearchResponse{Response: "服务端凭证响应回复收到：" + req.GetRequest()}, nil
}

func main() {
	// 根据服务端输入的证书文件和密钥构造 TLS 凭证
	c, err := credentials.NewServerTLSFromFile("./conf/tls/server-side/server.pem", "./conf/tls/server-side/server.key")
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
	}
	// 返回一个 ServerOption，用于设置服务器连接的凭据。
	// 用于 grpc.NewServer(opt ...ServerOption) 为 gRPC Server 设置连接选项
	s := grpc.NewServer(grpc.Creds(c))
	lis, err := net.Listen("tcp", ":8888") //创建 Listen，监听 TCP 端口
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
	}
	//将 SearchService（其包含需要被调用的服务端接口）注册到 gRPC Server 的内部注册中心。
	//这样可以在接受到请求时，通过内部的服务发现，发现该服务端接口并转接进行逻辑处理
	search.RegisterSearchServiceServer(s, &service{})

	//gRPC Server 开始 lis.Accept，直到 Stop 或 GracefulStop
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("gRPC Server 开始 lis.Accept err: %v", err)
		return
	}
}
