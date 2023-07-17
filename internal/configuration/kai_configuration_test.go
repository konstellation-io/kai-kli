package configuration_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kli/internal/configuration"
)

type KaiConfigurationTest struct {
	suite.Suite
}

func TestKaiConfigurationSuite(t *testing.T) {
	suite.Run(t, new(KaiConfigurationTest))
}

func (ch *KaiConfigurationTest) TestTokenIsValid_ExpectTrue() {
	// GIVEN
	token := configuration.Token{
		Date:             time.Now().UTC(),
		AccessToken:      "access_token",
		ExpiresIn:        3600,
		RefreshExpiresIn: 3600,
		RefreshToken:     "refresh_token",
		TokenType:        "Bearer",
	}

	// WHEN
	isValidToken := token.IsValid()

	// THEN
	ch.Require().True(isValidToken)
}

func (ch *KaiConfigurationTest) TestTokenIsValid_ExpectFalse() {
	// GIVEN
	token := configuration.Token{
		Date:             time.Now().AddDate(0, 0, -1).UTC(),
		AccessToken:      "access_token",
		ExpiresIn:        3600,
		RefreshExpiresIn: 3600,
		RefreshToken:     "refresh_token",
		TokenType:        "Bearer",
	}

	// WHEN
	isValidToken := token.IsValid()

	// THEN
	ch.Require().False(isValidToken)
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
	err := defaultConfig.AddServer(&server)

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
	err := defaultConfig.AddServer(&server)

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
	err := defaultConfig.AddServer(&server)

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
	err := defaultConfig.AddServer(&server)

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
	err := defaultConfig.AddServer(&server)

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

func (ch *KaiConfigurationTest) TestGetDefaultServer_GetDefaultServer_ExpectOk() {
	// GIVEN
	defaultConfig := configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{Name: "test1", URL: "http://test1.com", IsDefault: true},
			{Name: "test2", URL: "http://test2.com", IsDefault: false},
		},
	}

	// WHEN
	srv, err := defaultConfig.GetDefaultServer()

	// THEN
	ch.Require().NoError(err)
	ch.Require().Equal(&defaultConfig.Servers[0], srv)
}

func (ch *KaiConfigurationTest) TestGetDefaultServer_GetDefaultServerWhenNoServers_ExpectError() {
	// GIVEN
	defaultConfig := configuration.KaiConfiguration{
		Servers: []configuration.Server{},
	}

	// WHEN
	srv, err := defaultConfig.GetDefaultServer()

	// THEN
	ch.Require().Error(err)
	ch.Require().Nil(srv)
	ch.Require().ErrorIs(err, configuration.ErrNoServersFound)
}

func (ch *KaiConfigurationTest) TestUpdateServer_UpdateExistingServer_ExpectOK() {
	// GIVEN
	defaultConfig := configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{Name: "test1", URL: "http://test1.com", IsDefault: true, Realm: "realm1"},
		},
	}
	server := configuration.Server{
		Name:      "test1",
		URL:       "http://test1.com",
		IsDefault: false,
		Realm:     "realm2",
	}

	// WHEN
	err := defaultConfig.UpdateServer(&server)

	// THEN
	ch.Require().NoError(err)
	ch.Equal(1, len(defaultConfig.Servers))
	ch.Equal(configuration.Server{Name: "test1", URL: "http://test1.com", IsDefault: true, Realm: "realm2"},
		defaultConfig.Servers[0])
}

func (ch *KaiConfigurationTest) TestUpdateServer_UpdateExistingServer_WithDefault_ExpectOK() {
	// GIVEN
	defaultConfig := configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{Name: "test1", URL: "http://test1.com", IsDefault: false, Realm: "realm1"},
			{Name: "test2", URL: "http://test2.com", IsDefault: true, Realm: "realm2"},
		},
	}
	server := configuration.Server{
		Name:      "test1",
		URL:       "http://test1.com",
		IsDefault: true,
		Realm:     "realm2",
	}

	// WHEN
	err := defaultConfig.UpdateServer(&server)

	// THEN
	ch.Require().NoError(err)
	ch.Equal(2, len(defaultConfig.Servers))
	ch.Equal(configuration.Server{Name: "test1", URL: "http://test1.com", IsDefault: true, Realm: "realm2"},
		defaultConfig.Servers[0])
	ch.Equal(configuration.Server{Name: "test2", URL: "http://test2.com", IsDefault: false, Realm: "realm2"},
		defaultConfig.Servers[1])
}

func (ch *KaiConfigurationTest) TestUpdateServer_UpdateNonExistingServer_ExpectError() {
	// GIVEN
	defaultConfig := configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{Name: "test1", URL: "http://test1.com", IsDefault: true, Realm: "realm1"},
		},
	}
	server := configuration.Server{
		Name:      "test2",
		URL:       "http://test1.com",
		IsDefault: false,
		Realm:     "realm2",
	}

	// WHEN
	err := defaultConfig.UpdateServer(&server)

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, configuration.ErrServerNotFound)
}

func (ch *KaiConfigurationTest) TestDeleteServer_DeleteExistingServer_ExpectOK() {
	// GIVEN
	defaultConfig := configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{Name: "test1", URL: "http://test1.com", IsDefault: true, Realm: "realm1"},
			{Name: "test2", URL: "http://test2.com", IsDefault: false, Realm: "realm2"},
		},
	}

	// WHEN
	err := defaultConfig.DeleteServer("test1")

	// THEN
	ch.Require().NoError(err)
	ch.Equal(1, len(defaultConfig.Servers))
	ch.Equal(configuration.Server{Name: "test2", URL: "http://test2.com", IsDefault: true, Realm: "realm2"},
		defaultConfig.Servers[0])
}

func (ch *KaiConfigurationTest) TestDeleteServer_DeleteNonExistingServer_ExpectError() {
	// GIVEN
	defaultConfig := configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{Name: "test1", URL: "http://test1.com", IsDefault: true, Realm: "realm1"},
			{Name: "test2", URL: "http://test2.com", IsDefault: false, Realm: "realm2"},
		},
	}

	// WHEN
	err := defaultConfig.DeleteServer("server3")

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, configuration.ErrServerNotFound)
}

func (ch *KaiConfigurationTest) TestGetServer_GetExistingServer_ExpectOK() {
	// GIVEN
	defaultConfig := configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{Name: "test1", URL: "http://test1.com", IsDefault: true, Realm: "realm1"},
			{Name: "test2", URL: "http://test2.com", IsDefault: false, Realm: "realm2"},
		},
	}

	// WHEN
	srv, err := defaultConfig.GetServer("test1")

	// THEN
	ch.Require().NoError(err)
	ch.Require().NotNil(srv)
	ch.Equal(*srv, defaultConfig.Servers[0])
}

func (ch *KaiConfigurationTest) TestGetServer_GetNonExistingServer_ExpectError() {
	// GIVEN
	defaultConfig := configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{Name: "test1", URL: "http://test1.com", IsDefault: true, Realm: "realm1"},
			{Name: "test2", URL: "http://test2.com", IsDefault: false, Realm: "realm2"},
		},
	}

	// WHEN
	srv, err := defaultConfig.GetServer("server3")

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, configuration.ErrServerNotFound)
	ch.Require().Nil(srv)
}
