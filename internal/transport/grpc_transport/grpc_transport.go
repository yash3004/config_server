package grpc_transport

import (
	"context"
	"net"

	"github.com/yash3004/config_server/configurations"
	pb "github.com/yash3004/config_server/generated/protobuf/configpb"
	"github.com/yash3004/config_server/users"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedConfigServiceServer
	userManager   *users.UserManager
	configManager *configurations.ConfigManager
}

func NewServer(userManager *users.UserManager, configManager *configurations.ConfigManager) *Server {
	return &Server{
		userManager:   userManager,
		configManager: configManager,
	}
}

func StartGRPCServer(server *Server, address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterConfigServiceServer(s, server)
	return s.Serve(lis)
}

func (s *Server) AddConfig(ctx context.Context, req *pb.AddConfig) (*emptypb.Empty, error) {
	authenticated, err := s.userManager.AuthenticateUser(ctx, req.GetUserId(), req.GetPassword())
	if err != nil || !authenticated {
		return nil, err
	}

	err = s.configManager.AddConfig(ctx, req.GetUserId(), req.GetFilename(), int(req.GetFileType()), req.GetData())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateConfig(ctx context.Context, req *pb.UpdateConfig) (*emptypb.Empty, error) {
	authenticated, err := s.userManager.AuthenticateUser(ctx, req.GetUserId(), req.GetPassword())
	if err != nil || !authenticated {
		return nil, err
	}

	err = s.configManager.UpdateConfig(ctx, req.GetUserId(), req.GetFilename(), int(req.GetFileType()), req.GetData())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteConfig(ctx context.Context, req *pb.DeleteConfig) (*emptypb.Empty, error) {
	authenticated, err := s.userManager.AuthenticateUser(ctx, req.GetUserId(), req.GetPassword())
	if err != nil || !authenticated {
		return nil, err
	}

	err = s.configManager.DeleteConfig(ctx, req.GetUserId(), req.GetFilename())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) GetConfig(ctx context.Context, req *pb.GetConfig) (*pb.GetConfigResponse, error) {
	authenticated, err := s.userManager.AuthenticateUser(ctx, req.GetUserId(), req.GetPassword())
	if err != nil || !authenticated {
		return nil, err
	}

	data, fileType, err := s.configManager.GetConfig(ctx, req.GetUserId(), req.GetFilename())
	if err != nil {
		return nil, err
	}

	return &pb.GetConfigResponse{
		UserId:   req.GetUserId(),
		Filename: req.GetFilename(),
		FileType: pb.FileType(fileType),
		Data:     data,
	}, nil
}

func (s *Server) AddUser(ctx context.Context, req *pb.AddUser) (*emptypb.Empty, error) {
	err := s.userManager.AddUser(ctx, req.GetUserId(), req.GetEmail(), req.GetName(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUser) (*emptypb.Empty, error) {
	err := s.userManager.UpdateUser(ctx, req.GetUserId(), req.GetEmail(), req.GetName(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUser) (*emptypb.Empty, error) {
	err := s.userManager.DeleteUser(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}