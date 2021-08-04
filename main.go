// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/drone/drone-go/plugin/registry"
	"github.com/teryaev/drone-ecr-registry/plugin"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// spec provides the plugin settings.
type spec struct {
	Bind        string   `envconfig:"DRONE_BIND"`
	Debug       bool     `envconfig:"DRONE_DEBUG"`
	Secret      string   `envconfig:"DRONE_SECRET"`
	RegistryID  string   `envconfig:"ECR_REGISTRY_ID"`
	RegistryIDs []string `envconfig:"ECR_REGISTRY_IDS"`
}

//GetRegistries returns the registries in a backwards compatible manner
//such that ECR_REGISTRY_ID takes precendence over ECR_REGISTRY_IDS
func (s *spec) GetRegistries() (reg []string) {
	if s.RegistryID != "" {
		reg = []string{s.RegistryID}
	} else {
		reg = s.RegistryIDs
	}
	return
}

func main() {
	spec := new(spec)
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatal(err)
	}

	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if spec.Secret == "" {
		logrus.Fatalln("missing secret key")
	}
	if spec.Bind == "" {
		spec.Bind = ":3000"
	}

	cfg, err := external.LoadDefaultAWSConfig()

	if cfg.Region == "" {
		metaClient := ec2metadata.New(cfg)
		if region, err := metaClient.Region(context.Background()); err == nil {
			cfg.Region = region
			logrus.Infof("using region %s from ec2 metadata", cfg.Region)
		} else {
			logrus.Fatalf("failed to determine region: %s", err)
		}
	}

	if err != nil {
		logrus.Fatalln(err)
	}

	logrus.Debugf("Going to retrieve creds for %s ECR repo", spec.RegistryID)

	handler := registry.Handler(
		spec.Secret,
		plugin.New(
			ecr.New(cfg),
			spec.GetRegistries(),
		),
		logrus.StandardLogger(),
	)

	logrus.Infof("server listening on address %s", spec.Bind)

	http.Handle("/", handler)
	logrus.Fatal(http.ListenAndServe(spec.Bind, nil))
}
