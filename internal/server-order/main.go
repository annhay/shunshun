/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"shunshun/internal/pkg/initialization"
	"shunshun/internal/proto"
	"shunshun/internal/server-order/server"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	//初始化服务
	initialization.ServerInit()

	// 注册到 Consul
	consul := initialization.NewConsul("14.103.173.254:8500")
	kv := initialization.ConsulKV{
		Name:    "order-server",
		Tags:    []string{"order-server"},
		Address: "127.0.0.1",
		Port:    50053,
	}
	serviceID, err := consul.RegisterServer(kv)
	if err != nil {
		log.Printf("failed to register service: %v", err)
	}
	defer consul.DeregisterServer(serviceID)

	flag.Parse()
	lis, err := net.Listen("tcp", "127.0.0.1:50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterOrderServer(s, &server.Server{})
	log.Printf("server listening at %v", lis.Addr())

	// 启动服务
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 优雅关闭服务
	log.Println("Gracefully stopping server...")
	s.GracefulStop()
	log.Println("Server stopped gracefully")

	log.Println("Server exiting")
}
