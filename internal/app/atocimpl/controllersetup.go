package atocimpl

import (
	"context"

	"github.com/floffah/misecloud/api/atoc"
	"github.com/floffah/misecloud/internal/pkg/env"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RouteSetupServer struct {
	atoc.UnimplementedControllerSetupServer
}

func (r RouteSetupServer) GetIsSetup(ctx context.Context, empty *emptypb.Empty) (*atoc.GetIsSetupResponse, error) {
	isSetup := env.ControllerViper.GetBool(env.ControllerViperKeyIsSetup)
	return &atoc.GetIsSetupResponse{IsSetup: &isSetup}, nil
}

func (r RouteSetupServer) SetUp(ctx context.Context, req *atoc.SetupRequest) (*emptypb.Empty, error) {
	env.ControllerViper.Set(env.ControllerViperKeyIsSetup, true)
	env.ControllerViper.Set(env.ControllerViperKeyPort, req.GetPort())

	if err := env.ControllerViper.WriteConfig(); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
