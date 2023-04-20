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
	log.Println("Ping.")
	err := s.handler.Repository.Ping(ctx)
	if err != nil {
		return empty, status.Error(codes.Unavailable, "unavailable")
	}

	return empty, nil
}

func (s *GRPCServer) Stats(ctx context.Context, empty *emptypb.Empty) (*pb.StatsResponse, error) {
	stats := &pb.StatsResponse{}

	var addr string

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get(":authority")
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
		log.Println("nil from net.ParseIP: ", err)
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

func (s *GRPCServer) DeleteHandler(ctx context.Context, deletelist *pb.DeleteListRequest) (*emptypb.Empty, error) {
	fmt.Println("DeleteHandler.")
	empty := &emptypb.Empty{}
	fmt.Println("deletelist", deletelist)
	return empty, nil
}

func (s *GRPCServer) GetAllURL(ctx context.Context, empty *emptypb.Empty) (*pb.AllURLsResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unknown, "wrong metadata")
	}

	if len(md.Get("Cookies")) == 0 {
		return nil, status.Error(codes.Unknown, "wrong metadata")
	}

	cookie := md.Get("Cookies")[0]

	resp := &pb.AllURLsResponse{RespUrls: []*pb.PackReq{}}

	list, err := s.handler.Repository.GetAll(ctx, cookie)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}
	fmt.Println("list:", list)

	//resp.RespUrls = list

	//resp := &pb.AllURLsResponse{
	//	SortUrl:     list.ID,
	//	OriginalUrl: list.LongURL,
	//}

	fmt.Println("resp", resp)

	return resp, nil

}

func (s *GRPCServer) CreateShortURLHandler(ctx context.Context, body *pb.CreateShort) (*pb.CreateShortURLHandlerResponse, error) {
	log.Println("CreateShortURLHandler")

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unknown, "wrong metadata")
	}

	if len(md.Get("Cookies")) == 0 {
		return nil, status.Error(codes.Unknown, "wrong metadata")
	}
	cookie := md.Get("Cookies")[0]

	//генератор
	short := GenUlid()

	u, err := s.handler.Repository.SaveURL(ctx, short, body.Long, cookie)
	if err != nil {
		exist := &pb.CreateShortURLHandlerResponse{Url: s.cfg.BaseURL + "/" + u}
		return exist, status.Error(codes.AlreadyExists, "already exist")
	}

	newShort := s.cfg.BaseURL + "/" + short
	resp := &pb.CreateShortURLHandlerResponse{Url: newShort}
	return resp, nil
}

func (s *GRPCServer) GetShortURLByIDHandler(ctx context.Context) {

}

//
//func (s *GRPCServer) Batch() {
//
//}
//func (s *GRPCServer) CreateShortURLJSON() {
//
//}
