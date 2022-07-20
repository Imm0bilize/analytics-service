package v1

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client AnalyticsClient
}

func NewClient(host, port string) (*Client, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	grpcClient := NewAnalyticsClient(conn)
	return &Client{conn: conn, client: grpcClient}, nil
}

func (c *Client) CreateTask(ctx context.Context, taskID string) error {
	req := &NewTask{Id: taskID}
	_, err := c.client.CreateTask(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SetTimeStart(ctx context.Context, taskID, login, timeStart string, newTaskState *string) error {
	req := &TimeStart{User: &UsersTask{TaskId: taskID, Login: login}, TimeStart: timeStart, NewTaskState: newTaskState}
	_, err := c.client.SetTimeStart(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SetTimeEnd(ctx context.Context, taskID, login, timeEnd string, newTaskState *string) error {
	req := &TimeEnd{User: &UsersTask{TaskId: taskID, Login: login}, TimeEnd: timeEnd, NewTaskState: newTaskState}
	_, err := c.client.SetTimeEnd(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Shutdown() error {
	return c.conn.Close()
}
