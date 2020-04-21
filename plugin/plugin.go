// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/registry"
)

// New returns a new registry plugin.
func New(client *ecr.Client, registryID string) registry.Plugin {
	return &plugin{
		client:     client,
		registryID: registryID,
	}
}

type plugin struct {
	client     *ecr.Client
	registryID string
}

func (p *plugin) List(ctx context.Context, req *registry.Request) ([]*drone.Registry, error) {

	resp, err := p.client.GetAuthorizationTokenRequest(
		&ecr.GetAuthorizationTokenInput{
			RegistryIds: []string{p.registryID},
		},
	).Send(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("Couldn't retrieve auth token: %s", err.Error())
	}

	token := resp.AuthorizationData[0].AuthorizationToken
	decodedToken, err := base64.StdEncoding.DecodeString(*token)

	if err != nil {
		return nil, fmt.Errorf("Couldn't base64decode auth token: %s", err.Error())
	}

	creds := strings.Split(string(decodedToken), ":")

	credentials := []*drone.Registry{
		{
			Address:  *resp.AuthorizationData[0].ProxyEndpoint,
			Username: creds[0],
			Password: creds[1],
		},
	}

	return credentials, nil
}
