package service

import (
	"context"
	v1 "helloworld/api/helloworld/v1"
	"helloworld/internal/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	//main.LoadPromptDataAsync()
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	//if in.Name == "forLoop" {
	//	main.TestForLoop()
	//} else if in.Name == "query" {
	//	main.TestQuery()
	//}

	//g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	//if err != nil {
	//	return nil, err
	//}
	return &v1.HelloReply{Message: "Hello "}, nil
}
