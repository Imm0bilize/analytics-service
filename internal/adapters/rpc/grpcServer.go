package rpc

import (
	"analytic-service/internal/config"
	"analytic-service/internal/ports"
	v1 "analytic-service/pkg/commandsContract/v1"
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

type GrpcServer struct {
	v1.UnimplementedAnalyticsServer
	domain   ports.CommandDomain
	listener net.Listener
	log      logrus.FieldLogger
	notify   chan error
}

func New(cfg *config.Config, logger logrus.FieldLogger, domain ports.CommandDomain) (*GrpcServer, error) {
	listener, err := net.Listen(cfg.Grpc.Network, net.JoinHostPort(cfg.Grpc.Host, cfg.Grpc.Port))
	if err != nil {
		return nil, err
	}

	s := &GrpcServer{
		domain:   domain,
		notify:   make(chan error, 1),
		listener: listener,
		log:      logger,
	}
	return s, nil
}

func (s *GrpcServer) CreateTask(ctx context.Context, newTask *v1.NewTask) (*emptypb.Empty, error) {
	s.log.Infof("called")
	if err := s.domain.CreateTask(ctx, newTask.Id); err != nil {
		s.log.Errorf("error during execution: %s", err.Error())
		return nil, err
	}
	s.log.Infof("completed successfully")
	return &emptypb.Empty{}, nil
}

func (s *GrpcServer) SetTimeStart(ctx context.Context, timeStart *v1.TimeStart) (*emptypb.Empty, error) {
	s.log.Infof("called")
	if err := s.domain.SetTimeStart(ctx, timeStart.User.TaskId, timeStart.User.Login, timeStart.TimeStart, *timeStart.NewTaskState); err != nil {
		s.log.Errorf("error during execution: %s", err.Error())
		return nil, err
	}
	s.log.Infof("completed successfully")
	return &emptypb.Empty{}, nil
}

func (s *GrpcServer) SetTimeEnd(ctx context.Context, timeEnd *v1.TimeEnd) (*emptypb.Empty, error) {
	s.log.Infof("called")
	if err := s.domain.SetTimeEnd(ctx, timeEnd.User.TaskId, timeEnd.User.Login, timeEnd.TimeEnd, *timeEnd.NewTaskState); err != nil {
		s.log.Errorf("error during execution: %s", err.Error())
		return nil, err
	}
	s.log.Infof("completed successfully")
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
