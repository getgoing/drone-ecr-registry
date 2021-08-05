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
	"github.com/sirupsen/logrus"
)

// New returns a new registry plugin.
func New(client *ecr.Client, registries []string) registry.Plugin {
	return &plugin{
		client,
		registries,
	}
}

type plugin struct {
	client     *ecr.Client
	registries []string
}

func (p *plugin) List(ctx context.Context, req *registry.Request) ([]*drone.Registry, error) {

	resp, err := p.client.GetAuthorizationTokenRequest(
		&ecr.GetAuthorizationTokenInput{
			RegistryIds: p.registries,
		},
	).Send(context.TODO())

	if err != nil {
		result := fmt.Errorf("couldn't retrieve auth token from registries %#v: %w", p.registries, err)
		logrus.Error(result)
		return nil, result
	}

	credentials := make([]*drone.Registry, 0)
	for _, authData := range resp.AuthorizationData {
		token := authData.AuthorizationToken
		if decodedToken, err := base64.StdEncoding.DecodeString(*token); err != nil {
			result := fmt.Errorf("couldn't decode auth token for %s: %w", *authData.ProxyEndpoint, err)
			logrus.Error(result)
			return nil, result
		} else {
			creds := strings.Split(string(decodedToken), ":")
			credentials = append(credentials, &drone.Registry{
				Address:  *authData.ProxyEndpoint,
				Username: creds[0],
				Password: creds[1],
			})
		}
	}

	return credentials, nil
}
