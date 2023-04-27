package handler

import (
	"context"
	"net"

	"github.com/22Fariz22/shorturl/internal/config"
	"github.com/22Fariz22/shorturl/internal/entity"
	"github.com/22Fariz22/shorturl/pkg/logger"
	pb "github.com/22Fariz22/shorturl/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	pb.UnimplementedServicesServer
	l       logger.Interface
	cfg     config.Config
	handler *Handler
}

func NewGRPCServer(l logger.Interface, cfg config.Config, handler *Handler) *GRPCServer {
	return &GRPCServer{
		UnimplementedServicesServer: pb.UnimplementedServicesServer{},
		l:                           l,
		cfg:                         cfg,
		handler:                     handler,
	}
}

func (s *GRPCServer) Ping(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	err := s.handler.Repository.Ping(ctx, s.l)
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
		s.l.Info("err net.SplitHostPort: ", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	// парсим ip
	ip := net.ParseIP(ipStr)
	if ip == nil {
		s.l.Info("nil from net.ParseIP: ", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	_, ipnet, err := net.ParseCIDR(s.cfg.TrustedSubnet)
	if err != nil {
		s.l.Info("err net.ParseCIDR: ", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	if ipnet.Contains(ip) {
		urls, users, err := s.handler.Repository.Stats(ctx, s.l)
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
	empty := &emptypb.Empty{}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unknown, "wrong metadata")
	}

	if len(md.Get("Cookies")) == 0 {
		return nil, status.Error(codes.Unknown, "wrong metadata")
	}

	cookie := md.Get("Cookies")[0]

	arr := make([]string, len(deletelist.DeleteList))

	for _, v := range deletelist.DeleteList {
		arr = append(arr, v.OneString)
	}

	s.handler.Workers.AddJob(ctx, s.l, arr, cookie)

	return empty, nil
}

func (s *GRPCServer) GetAllURL(ctx context.Context, empty *emptypb.Empty) (*pb.AllURLRequestList, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unknown, "wrong metadata")
	}

	if len(md.Get("Cookies")) == 0 {
		return nil, status.Error(codes.Unknown, "wrong metadata")
	}

	cookie := md.Get("Cookies")[0]

	repoAnswer, err := s.handler.Repository.GetAll(ctx, s.l, cookie)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	//загатовка для request
	batchListReq := make([]*pb.AllURLRequest, 0, len(repoAnswer))

	for _, mp := range repoAnswer {
		for k, v := range mp {
			temp := &pb.AllURLRequest{
				ShortUrl:    s.cfg.BaseURL + "/" + k,
				OriginalUrl: v,
			}
			//добавялем в загатовку для request
			batchListReq = append(batchListReq, temp)
		}
	}

	//для request
	resp := &pb.AllURLRequestList{AllUrls: batchListReq}

	return resp, nil
}

func (s *GRPCServer) CreateShortURLHandler(ctx context.Context, body *pb.CreateShort) (*pb.CreateShortURLHandlerResponse, error) {
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

	u, err := s.handler.Repository.SaveURL(ctx, s.l, short, body.Long, cookie)
	if err != nil {
		exist := &pb.CreateShortURLHandlerResponse{Url: s.cfg.BaseURL + "/" + u}
		return exist, status.Error(codes.AlreadyExists, "already exist")
	}

	newShort := s.cfg.BaseURL + "/" + short
	resp := &pb.CreateShortURLHandlerResponse{Url: newShort}
	return resp, nil
}

func (s *GRPCServer) GetShortURLByIDHandler(ctx context.Context, param *pb.IDParam) (*pb.OneString, error) {
	url, ok := s.handler.Repository.GetURL(ctx, s.l, param.Id)
	if !ok {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	header := metadata.Pairs("Location", url.LongURL)
	grpc.SendHeader(ctx, header)

	resp := &pb.OneString{OneString: url.LongURL}

	return resp, nil
}

func (s *GRPCServer) Batch(ctx context.Context, packReq *pb.BatchListResp) (*pb.BatchListReq, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unknown, "wrong metadata")
	}

	if len(md.Get("Cookies")) == 0 {
		return nil, status.Error(codes.Unknown, "wrong metadata")
	}

	cookie := md.Get("Cookies")[0]

	//загатовка message response
	arrPackReq := make([]entity.PackReq, 0, len(packReq.GetBatchListResp()))

	//загатовка message  request
	arrBatchReq := make([]*pb.BatchReq, 0, len(packReq.GetBatchListResp()))

	//добавляем в массив для отправки в репу и другой массив для response
	for _, v := range packReq.GetBatchListResp() {
		//генератор
		short := GenUlid()

		//заполняем заготовку для репы
		preq := entity.PackReq{
			CorrelationID: v.CorrelationId,
			OriginalURL:   v.OriginalUrl,
			ShortURL:      short,
		}
		arrPackReq = append(arrPackReq, preq)

		//заполняем заготовку для request
		batchReq := &pb.BatchReq{
			CorrelationId: v.CorrelationId,
			ShortUrl:      short,
		}
		arrBatchReq = append(arrBatchReq, batchReq)

	}

	//отравляем заполненную загатовку в репу
	err := s.handler.Repository.RepoBatch(ctx, s.l, cookie, arrPackReq)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	//окончательная загатовка для request
	bResp := &pb.BatchListReq{BatchListResult: arrBatchReq}

	return bResp, nil
}

func (s *GRPCServer) CreateShortURLJSON(ctx context.Context, res *pb.ReqURL) (*pb.CreateShortURLJSONResponse, error) {
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

	r, err := s.handler.Repository.SaveURL(ctx, s.l, short, res.Url, cookie) //если есть такой,то вернуть шорт и конфликт статус
	if err != nil {
		if r != "" {
			exist := s.cfg.BaseURL + "/" + r
			return &pb.CreateShortURLJSONResponse{Result: exist}, nil
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	resp := s.cfg.BaseURL + "/" + short

	return &pb.CreateShortURLJSONResponse{Result: resp}, nil
}
