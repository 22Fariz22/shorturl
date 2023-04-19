package handler

import (
	"context"
	"fmt"
	"github.com/22Fariz22/shorturl/internal/config"
	pb "github.com/22Fariz22/shorturl/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

type GRPCServer struct {
	pb.UnimplementedServicesServer
	cfg     config.Config
	handler *Handler
}

func NewGRPCServer(cfg config.Config, handler *Handler) *GRPCServer {
	return &GRPCServer{
		UnimplementedServicesServer: pb.UnimplementedServicesServer{},
		cfg:                         cfg,
		handler:                     handler,
	}
}

func (s *GRPCServer) Ping(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	empty = &emptypb.Empty{}

	err := s.handler.Repository.Ping(ctx)
	if err != nil {
		return empty, status.Error(codes.Unavailable, "unavailable")
	}

	return empty, nil
}

func (s *GRPCServer) Stats(ctx context.Context, empty *emptypb.Empty) (*pb.StatsResponse, error) {
	stats := &pb.StatsResponse{}

	//addr := r.RemoteAddr
	var addr string

	md, ok := metadata.FromIncomingContext(ctx)
	fmt.Println("md: ", md)
	if ok {
		values := md.Get(":authority")
		fmt.Println("addr", values)
		if len(values) > 0 {
			// ключ содержит слайс строк, получаем первую строку
			addr = values[0]
		}
	}

	ipStr, _, err := net.SplitHostPort(addr)
	if err != nil {
		log.Println("err net.SplitHostPort: ", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	// парсим ip
	ip := net.ParseIP(ipStr)
	if ip == nil {
		log.Println("err net.ParseIP: ", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	_, ipnet, err := net.ParseCIDR(s.cfg.TrustedSubnet)
	if err != nil {
		log.Println("err net.ParseCIDR: ", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	if ipnet.Contains(ip) {
		urls, users, err := s.handler.Repository.Stats(ctx)
		if err != nil {
			return nil, status.Error(codes.Internal, "internal server error")
		}
		stats.Urls = int32(urls)
		stats.Users = int32(users)
	} else {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return stats, nil
}

func DeleteHandler() {

}
func GetAllURL() {

}
func CreateShortURLHandler() {

}
func GetShortURLByIDHandler() {

}
func Batch() {

}
func CreateShortURLJSON() {

}
