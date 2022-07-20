package integration_test

import (
	v1 "analytic-service/pkg/commandsContract/v1"
	"context"
	"github.com/stretchr/testify/require"
)

func (i *integrationTestSuite) TestCommandPositive() {
	client, err := v1.NewClient(serviceUrl, serviceGrpcPort)
	if err != nil {
		i.T().Fatal(err)
	}
	defer client.Shutdown()
	ctx := context.Background()

	err = client.CreateTask(ctx, "testing")
	require.Nil(i.T(), err)

	newState := "processing"
	err = client.SetTimeStart(ctx, "testing", "test1", "2022-07-18 18:30", &newState)
	require.Nil(i.T(), err)

	newState = "accepted"
	err = client.SetTimeEnd(ctx, "testing", "test1", "2022-07-18 19:00", &newState)
	require.Nil(i.T(), err)
}

func (i *integrationTestSuite) TestCommandNegative() {
	client, err := v1.NewClient(serviceUrl, serviceGrpcPort)
	if err != nil {
		i.T().Fatal(err)
	}
	defer client.Shutdown()
	ctx := context.Background()

	newState := "processing"
	err = client.SetTimeStart(ctx, "incorrect_id", "test1", "2022-07-18 18:30", &newState)
	require.NotNil(i.T(), err)
}
