package main

import (
	"context"
	"flag"
	"fmt"
	st "go-web/stream"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

var (
	port = flag.Int("port", 40111, "The server port")
)

// 必须嵌入UnimplementedUserServiceServer
type server struct {
	st.UnimplementedGreeterServer
}

// Single 实现Single-一元简单RPC方法
func (s *server) Single(ctx context.Context, req *st.StreamReqData) (resp *st.StreamResData, err error) {
	log.Printf("收到摩斯密码请求: %v", req.String())
	return &st.StreamResData{Message: "哈喽,我在 <一元RPC> 收到了你的摩斯密码: \"" + req.MorseCode + "\"\t"}, nil
}

/*
GetStream 服务端流式RPC，Server是Stream，Client为普通RPC请求
客户端发送一次普通的RPC请求，服务端通过流式响应多次发送数据集
*/
func (s *server) GetStream(req *st.StreamReqData, gss st.Greeter_GetStreamServer) error {

	log.Printf("收到摩斯密码请求: %v", req.String())

	// 具体返回多少个response根据业务逻辑调整
	for n := 0; n < 6; n++ {
		// 通过 send 方法不断推送数据
		err := gss.Send(&st.StreamResData{Message: "哈喽,我在 <服务端流> 收到了你的摩斯密码: \"" + req.GetMorseCode() + "\""})
		if err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	// 返回nil表示已经完成响应
	return nil
}

/*
PutStream 客户端流式RPC，单向流
客户端通过流式多次发送RPC请求给服务端，服务端发送一次普通的RPC请求给客户端
*/
func (s *server) PutStream(pss st.Greeter_PutStreamServer) error {
	//for循环接收客户端发送的消息
	for {
		// 通过 Recv() 不断获取客户端 send()推送的消息
		req, err := pss.Recv()
		// err == io.EOF表示已经获取全部数据
		if err == io.EOF {
			// SendAndClose 返回并关闭连接，在客户端发送完毕后服务端即可返回响应
			return pss.SendAndClose(&st.StreamResData{Message: "再见,我在 <客户端流> 收完了你全部的摩斯密码"})
		}
		if err != nil {
			log.Fatalf("---PutStream--- 出错了 Error: %s", err)
			return err
		}
		log.Printf("---PutStream--- 收到摩斯密码请求: %s\n", req.GetMorseCode())
		time.Sleep(time.Second)
	}
	//返回nil表示已经完成响应
	return nil
}

/*
AllStream 双向流，由客户端发起流式的RPC方法请求，服务端以同样的流式RPC方法响应请求
首个请求一定是client发起，具体交互方法（谁先谁后，一次发多少，响应多少，什么时候关闭）根据程序编写方式来确定（可以结合协程）
1. 建立连接 获取client
// 2. 通过client调用方法获取stream
// 3. 开两个goroutine（使用 chan 传递数据） 分别用于Recv()和Send()
// 3.1 一直Recv()到err==io.EOF(即客户端关闭stream)
// 3.2 Send()则自己控制什么时候Close 服务端stream没有close 只要跳出循环就算close了。 具体见https://github.com/grpc/grpc-go/issues/444
*/
func (s *server) AllStream(stream st.Greeter_AllStreamServer) error {
	var (
		wg    sync.WaitGroup //任务编排
		msgCh = make(chan *string)
	)

	wg.Add(1)
	go func() {
		n := 0
		defer wg.Done()
		for v := range msgCh {
			err := stream.Send(&st.StreamResData{Message: "收到摩斯密码请求：" + *v})
			if err != nil {
				log.Printf("---AllStream--- 回复出错了 Error: %s\n", err)
				continue
			}
			n++
			log.Printf("---AllStream--- N: %d\n", n)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("recv error :%v", err)
			}
			log.Printf("--AllStream--- 收到摩斯密码%s\n", req.MorseCode)
			msgCh <- &req.MorseCode
		}
		close(msgCh)
	}()

	wg.Wait() //等待任务结束

	return nil
}

func main() {
	flag.Parse()
	// 设置监听端口
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcSer := grpc.NewServer()

	st.RegisterGreeterServer(grpcSer, &server{})
	log.Printf("server listening at %v \n", lis.Addr())

	if err := grpcSer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
