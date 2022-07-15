package rpc

import (
	"analytic-service/internal/config"
	"analytic-service/internal/ports"
	v1 "analytic-service/pkg/interpretedProto/v1"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

type GrpcServer struct {
	v1.UnimplementedAnalyticsServer
	domain   ports.ManagementServerDomain
	notify   chan error
	listener net.Listener
}

func New(cfg *config.Config, domain ports.ManagementServerDomain) (*GrpcServer, error) {
	listener, err := net.Listen(cfg.Grpc.Network, net.JoinHostPort(cfg.Grpc.Host, cfg.Grpc.Port))
	if err != nil {
		return nil, err
	}

	s := &GrpcServer{
		domain:   domain,
		notify:   make(chan error, 1),
		listener: listener,
	}
	return s, nil
}

func (s *GrpcServer) CreateTask(ctx context.Context, newTask *v1.NewTask) (*emptypb.Empty, error) {
	if err := s.domain.CreateTask(ctx, newTask.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *GrpcServer) UpdateTasksState(ctx context.Context, newState *v1.NewTasksState) (*emptypb.Empty, error) {
	if err := s.domain.UpdateTasksState(ctx, newState.Id, newState.State); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *GrpcServer) SetTimeStart(ctx context.Context, timeStart *v1.TimeStart) (*emptypb.Empty, error) {
	if err := s.domain.SetTimeStart(ctx, timeStart.User.TaskId, timeStart.User.Login, timeStart.TimeStart); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *GrpcServer) SetTimeEnd(ctx context.Context, timeEnd *v1.TimeEnd) (*emptypb.Empty, error) {
	if err := s.domain.SetTimeStart(ctx, timeEnd.User.TaskId, timeEnd.User.Login, timeEnd.TimeEnd); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *GrpcServer) Notify() <-chan error {
	return s.notify
}

func (s *GrpcServer) Run() {
	go func() {
		var opts []grpc.ServerOption
		grpcServer := grpc.NewServer(opts...)
		v1.RegisterAnalyticsServer(grpcServer, s)

		if err := grpcServer.Serve(s.listener); err != nil {
			s.notify <- err
			close(s.notify)
		}
	}()
}

func (s *GrpcServer) Shutdown() error {
	return s.listener.Close()
}
