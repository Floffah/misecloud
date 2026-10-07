package atocimpl

import (
	"context"

	"github.com/floffah/misecloud/api/atoc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RouteHealthServer struct {
	atoc.UnimplementedHealthServer
}

func (r RouteHealthServer) Check(ctx context.Context, empty *emptypb.Empty) (*atoc.HealthCheckResponse, error) {
	status := atoc.ServingStatus_SERVING

	return &atoc.HealthCheckResponse{
		Status: &status,
	}, nil
}
