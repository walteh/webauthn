package dynamodb

import (
	"context"

	"github.com/walteh/buildrc/integration"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var _ integration.ContainerImage = (*DockerImage)(nil)

type DockerImage struct {
	active *integration.ContainerStore
}

func (me *DockerImage) OnStart(z *integration.ContainerStore) {
	me.active = z
}

func (me *DockerImage) Tag() string {
	return "amazon/dynamodb-local:latest"
}

func (me *DockerImage) HTTPPort() int {
	return 8000
}

func (me *DockerImage) HTTPSPort() int {
	return 8000
}

func (me *DockerImage) EnvVars() []string {
	return []string{}
}

func (me *DockerImage) Ping(ctx context.Context) error {
	c, err := me.NewClient()
	if err != nil {
		return err
	}
	_, err = c.DescribeLimits(ctx, &dynamodb.DescribeLimitsInput{})
	return err
}
