package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/floffah/misecloud/api/atoc"
	"github.com/floffah/misecloud/internal/app/atocimpl"
	"github.com/floffah/misecloud/internal/pkg/env"
	"github.com/floffah/misecloud/internal/pkg/tailscale"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	err := env.InitControllerViper()
	if err != nil {
		panic(fmt.Errorf("could not initialise environment: %s", err.Error()))
	}

	tsDaemonStatus, err := tailscale.LocalClient.Status(ctx)
	if err != nil {
		panic(fmt.Errorf("could not get Tailscale daemon status: %s", err.Error()))
	}

	if len(tsDaemonStatus.TailscaleIPs) == 0 {
		panic(fmt.Errorf("no Tailscale IPs found, is Tailscale running"))
	}

	if !env.ControllerViper.GetBool(env.ControllerViperKeyIsSetup) {
		hostname := strings.TrimSuffix(tsDaemonStatus.Self.DNSName, ".")
		fmt.Printf("Controller is not set up. Please run `misecl controller setup %s`\n", hostname)
	}

	address := net.JoinHostPort(tsDaemonStatus.TailscaleIPs[0].String(), env.ControllerViper.GetString(env.ControllerViperKeyPort))
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	atoc.RegisterHealthServer(grpcServer, atocimpl.RouteHealthServer{})
	atoc.RegisterControllerSetupServer(grpcServer, atocimpl.RouteSetupServer{})

	fmt.Printf("Listening at %s\n", address)
	grpcServer.Serve(lis)
}
