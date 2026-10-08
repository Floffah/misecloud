package tailscale

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"tailscale.com/client/local"
	tsremote "tailscale.com/client/tailscale/v2"
)

var ErrNotAuthed = errors.New("tailscale remote api not authenticated in keyring")

var LocalClient = &local.Client{}

var remoteClient *tsremote.Client

func EnsureControllerTag(ctx context.Context, client *tsremote.Client, deviceID string) error {
	if deviceID == "" {
		return errors.New("controller did not provide a Tailscale device ID; update the controller before setup")
	}
	device, err := client.Devices().Get(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("get controller device: %w", err)
	}
	if slices.Contains(device.Tags, "tag:misecl-controller") {
		return nil
	}
	if err := client.Devices().SetTags(ctx, deviceID, append(device.Tags, "tag:misecl-controller")); err != nil {
		return fmt.Errorf("tag controller device (define tag:misecl-controller and grant OAuth Devices Core write access for this tag): %w", err)
	}
	return nil
}

func GetRemoteApi(ctx context.Context) (*tsremote.Client, error) {
	if remoteClient == nil {
		status, err := LocalClient.Status(ctx)
		if err != nil {
			return nil, err
		}

		tscreds, err := GetOAuthCredentials(string(status.CurrentTailnet.StableID))
		if err != nil {
			return nil, err
		}
		if !tscreds.IsAuthed {
			return nil, ErrNotAuthed
		}

		remoteClient = &tsremote.Client{
			Tailnet: string(status.CurrentTailnet.StableID),
			Auth: &tsremote.OAuth{
				ClientID:     tscreds.ClientID,
				ClientSecret: tscreds.ClientSecret,
			},
		}
	}

	return remoteClient, nil
}
