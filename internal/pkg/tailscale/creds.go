package tailscale

import (
	"encoding/json"
	"errors"

	"github.com/zalando/go-keyring"
)

const keyringOauthName = "com.github.floffah.misecloud.misecl.tailscale.oauth"

type OAuthCredentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type OAuthCredentialsResponse struct {
	OAuthCredentials
	IsAuthed bool
}

func GetOAuthCredentials(profileName string) (OAuthCredentialsResponse, error) {
	strinfo, err := keyring.Get(keyringOauthName, profileName)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return OAuthCredentialsResponse{IsAuthed: false}, nil
		}
		return OAuthCredentialsResponse{IsAuthed: false}, err
	}

	var creds OAuthCredentials
	err = json.Unmarshal([]byte(strinfo), &creds)
	if err != nil {
		return OAuthCredentialsResponse{IsAuthed: false}, err
	}

	return OAuthCredentialsResponse{
		OAuthCredentials: creds,
		IsAuthed:         true,
	}, nil
}

func SetOAuthCredentials(profileName string, creds OAuthCredentials) error {
	strinfo, err := json.Marshal(creds)
	if err != nil {
		return err
	}

	err = keyring.Set(keyringOauthName, profileName, string(strinfo))
	if err != nil {
		return err
	}

	return nil
}
