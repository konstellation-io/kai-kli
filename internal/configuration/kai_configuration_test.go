package configuration_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kli/internal/configuration"
)

type KaiConfigurationTest struct {
	suite.Suite
}

func TestKaiConfigurationSuite(t *testing.T) {
	suite.Run(t, new(KaiConfigurationTest))
}

func (ch *KaiConfigurationTest) TestAddServer_AddFirstServer_ExpectOK() {
	// GIVEN
	defaultConfig := configuration.KaiConfiguration{}
	server := configuration.Server{
		Name:      "test",
		URL:       "http://test.com",
		IsDefault: false,
	}

	// WHEN
	err := defaultConfig.AddServer(server)

	// THEN
	ch.Require().NoError(err)
	ch.Equal(1, len(defaultConfig.Servers))
	ch.Equal(configuration.Server{Name: "test", URL: "http://test.com", IsDefault: true}, defaultConfig.Servers[0])
}

func (ch *KaiConfigurationTest) TestAddServer_AddNewServerToExistingList_ExpectOK() {
	// GIVEN
	defaultConfig := configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{Name: "test1", URL: "http://test1.com", IsDefault: true},
		},
	}
	server := configuration.Server{
		Name:      "test2",
		URL:       "http://test2.com",
		IsDefault: false,
	}

	// WHEN
	err := defaultConfig.AddServer(server)

	// THEN
	ch.Require().NoError(err)
	ch.Equal(2, len(defaultConfig.Servers))
	ch.Equal(configuration.Server{Name: "test1", URL: "http://test1.com", IsDefault: true}, defaultConfig.Servers[0])
	ch.Equal(configuration.Server{Name: "test2", URL: "http://test2.com", IsDefault: false}, defaultConfig.Servers[1])
}

func (ch *KaiConfigurationTest) TestAddServer_AddNewDefaultServerToExistingList_ExpectOK() {
	// GIVEN
	defaultConfig := configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{Name: "test1", URL: "http://test1.com", IsDefault: true},
		},
	}
	server := configuration.Server{
		Name:      "test2",
		URL:       "http://test2.com",
		IsDefault: true,
	}

	// WHEN
	err := defaultConfig.AddServer(server)

	// THEN
	ch.Require().NoError(err)
	ch.Equal(2, len(defaultConfig.Servers))
	ch.Equal(configuration.Server{Name: "test1", URL: "http://test1.com", IsDefault: false}, defaultConfig.Servers[0])
	ch.Equal(configuration.Server{Name: "test2", URL: "http://test2.com", IsDefault: true}, defaultConfig.Servers[1])
}

func (ch *KaiConfigurationTest) TestAddServer_AddDuplicatedServerName_ExpectError() {
	// GIVEN
	defaultConfig := configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{Name: "test1", URL: "http://test1.com", IsDefault: true},
		},
	}
	server := configuration.Server{
		Name:      "test1",
		URL:       "http://test2.com",
		IsDefault: true,
	}

	// WHEN
	err := defaultConfig.AddServer(server)

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, configuration.ErrDuplicatedServerName)
	ch.Equal(1, len(defaultConfig.Servers))
}

func (ch *KaiConfigurationTest) TestAddServer_AddDuplicatedServerUrl_ExpectError() {
	// GIVEN
	defaultConfig := configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{Name: "test1", URL: "http://test1.com", IsDefault: true},
		},
	}
	server := configuration.Server{
		Name:      "test2",
		URL:       "http://test1.com",
		IsDefault: true,
	}

	// WHEN
	err := defaultConfig.AddServer(server)

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, configuration.ErrDuplicatedServerURL)
	ch.Equal(1, len(defaultConfig.Servers))
}

func (ch *KaiConfigurationTest) TestSetDefaultServer_ChangeDefaultServerToExistingServer_ExpectOK() {
	// GIVEN
	defaultConfig := configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{Name: "test1", URL: "http://test1.com", IsDefault: true},
			{Name: "test2", URL: "http://test2.com", IsDefault: false},
		},
	}

	// WHEN
	err := defaultConfig.SetDefaultServer("test2")

	// THEN
	ch.Require().NoError(err)
	ch.Equal(configuration.Server{Name: "test1", URL: "http://test1.com", IsDefault: false}, defaultConfig.Servers[0])
	ch.Equal(configuration.Server{Name: "test2", URL: "http://test2.com", IsDefault: true}, defaultConfig.Servers[1])
}

func (ch *KaiConfigurationTest) TestSetDefaultServer_ChangeDefaultServerToNonExistingServer_ExpectError() {
	// GIVEN
	defaultConfig := configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{Name: "test1", URL: "http://test1.com", IsDefault: true},
			{Name: "test2", URL: "http://test2.com", IsDefault: false},
		},
	}

	// WHEN
	err := defaultConfig.SetDefaultServer("test3")

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, configuration.ErrServerNotFound)
}
