package main

import (
	"os"
	"testing"

	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SpecTestSuite struct {
	suite.Suite
}

func TestSpecTestSuite(t *testing.T) {
	suite.Run(t, new(SpecTestSuite))
}

func (suite *SpecTestSuite) TestRegistries() {
	t := suite.T()
	s := new(spec)
	os.Setenv("ECR_REGISTRY_IDS", "reg1")
	err := envconfig.Process("", s)
	assert.NoError(t, err)
	assert.Equal(t, []string{"reg1"}, s.RegistryIDs)
	os.Setenv("ECR_REGISTRY_IDS", "rega,regb")
	s = new(spec)
	err = envconfig.Process("", s)
	assert.NoError(t, err)
	assert.Equal(t, []string{"rega", "regb"}, s.RegistryIDs)
}
