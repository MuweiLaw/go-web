package main

import (
	"context"
	"flag"
	"fmt"
	st "go-web/stream"
	"io"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "hello,world!!"
)

var (
	//addr = flag.String("addr", "139.159.191.200:40111", "the address to connect to")
	addr = flag.String("addr", "localhost:40111", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// 一定要记得关闭链接
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("关闭链接错误: %v", err)
		}
	}(conn)

	//定义请求体
	req := st.StreamReqData{MorseCode: *name}

	// --------------------------------------  一元GRPC  ------------------------------------------------

	// 实例化客户端
	client := st.NewGreeterClient(conn)
	// 直接发起请求
	resp, err := client.Single(context.Background(), &req)
	log.Printf("---Single--- <直接发出>发出摩斯密码 ===>> \"%s\" \n", req.GetMorseCode())
	if err != nil {
		log.Fatalf("---Single--- <直接发出>Error:\n %v \n", err)
	}
	log.Printf("---Single--- <直接发出>回应消息 ===>> %s \n\n", resp.GetMessage())

	// 联系服务器并打印出它的响应
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	// 超时关闭客户端
	defer cancel()

	req.MorseCode = "在嘛, 我在 <一元RPC> #定义上下文超时# 我的头又秃了一点"
	resp, err = client.Single(ctx, &req)
	log.Printf("---Single--- <定义上下文超时>发出摩斯密码 ===>> \"%s\" \n", req.GetMorseCode())

	if err != nil {
		log.Fatalf("---Single--- <定义上下文超时>Error:|n %v \n", err)
	}
	log.Printf("---Single--- <定义上下文超时>回应消息 ===>> %s \n\n", resp.GetMessage())

	// ---------------------------------------------------  进阶,流GRPC  ----------------------------------------------

	req.MorseCode = "在嘛, 服务器流, 我的头又秃了一点"
	getStreamClient, err := client.GetStream(context.Background(), &req)
	if err != nil {
		log.Fatalf("---GetStream--- GetStream Error: \n%v", err)
	}
	_ = printGetStream(getStreamClient, &req)
	log.Print("---GetStream--- 服务器流结束\n\n")

	req.MorseCode = "在嘛, 客户端流, 我的头又秃了一点"
	putStreamClient, err := client.PutStream(context.Background())
	if err != nil {
		log.Fatalf("---PutStream--- PutStream Error: \n%v", err)
	}
	_ = printPutStream(putStreamClient, &req)
	log.Print("---GetStream--- 客户端流结束\n\n")

	req.MorseCode = "在嘛, 双向流, 我的头又秃了一点"
	allStreamClient, err := client.AllStream(context.Background())
	if err != nil {
		log.Fatalf("---AllStream--- AllStream Error: \n%v", err)
	}
	_ = printAllStream(allStreamClient, &req)
}

/*
1. 建立连接 获取client
2. 通过 client 获取stream
3. for循环中通过stream.Recv()依次获取服务端推送的消息
4. err==io.EOF则表示服务端关闭stream了
*/
func printGetStream(client st.Greeter_GetStreamClient, req *st.StreamReqData) error {

	// 发送一个请求消息
	err := client.CloseSend()
	log.Printf("---GetStream--- 发出摩斯密码: %s", req.GetMorseCode())
	if err != nil {
		log.Fatalf("---GetStream--- 发出摩斯密码 Error \n%s", err)
		return err
	}

	// for循环获取服务端推送的消息
	for {
		// 通过 Recv() 不断获取服务端send()推送的消息
		resp, err := client.Recv()
		// err==io.EOF则表示服务端关闭stream了
		if err == io.EOF {
			log.Printf("---GetStream--- 服务端关闭stream EOF: %s", err)
			break
		}
		if err != nil {
			log.Fatalf("---GetStream--- 接受回应消息 Error \n%s", err)
			return err
		}
		log.Printf("---GetStream--- 回应消息 ===>> %s", resp.GetMessage())
	}
	return nil
}

func printPutStream(client st.Greeter_PutStreamClient, req *st.StreamReqData) error {

	for i := 0; i < 6; i++ {
		//通过 SendMsg 方法不断推送数据到服务端
		err := client.SendMsg(req)
		log.Printf("---PutStream--- 发出摩斯密码: %s", req.GetMorseCode())
		if err != nil {
			log.Fatalf("---PutStream--- 发出摩斯密码 Error \n%s", err)
			return err
		}
	}

	// 发送完成后通过stream.CloseAndRecv() 关闭stream并接收服务端返回结果
	// 服务端则根据err==io.EOF来判断client是否关闭stream
	resp, err := client.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("---PutStream--- 回应消息 ===>> %s", resp.GetMessage())
	return nil
}

/*
1. 建立连接 获取client
2. 通过client获取stream
3. 开两个goroutine 分别用于Recv()和Send()
3.1 一直Recv()到err==io.EOF(即服务端关闭stream)
3.2 Send()则由自己控制
4. 发送完毕调用 stream.CloseSend()关闭stream 必须调用关闭 否则Server会一直尝试接收数据 一直报错...
*/
func printAllStream(client st.Greeter_AllStreamClient, r *st.StreamReqData) error {
	var wg sync.WaitGroup

	// 开两个goroutine 分别用于Recv()和Send()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			resp, err := client.Recv()
			if err == io.EOF {
				fmt.Println("---AllStream--- Server Closed")
				break
			}
			if err != nil {
				continue
			}
			log.Printf("---AllStream--- Recv Message: %s", resp.Message)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for n := 0; n < 6; n++ {
			err := client.Send(r)
			if err != nil {
				log.Printf("---AllStream--- Send Error:%v\n", err)
			}
			log.Printf("---AllStream--- Send Succeed:%v\n", r)
			time.Sleep(time.Second)
		}

		// 发送完毕关闭stream
		err := client.CloseSend()
		if err != nil {
			log.Printf("---AllStream--- Close Send Error:%v\n", err)
			return
		}
	}()

	wg.Wait()
	return nil
}
