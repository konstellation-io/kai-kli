package api_test

import (
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/testhelpers"
)

func TestNewKaiClient(t *testing.T) {
	d := testhelpers.SetupConfigDir(t)
	defer testhelpers.CleanConfigDir(t, d)

	viper.Set(config.RequestTimeoutKey, 2*time.Minute)
	viper.Set(config.DebugKey, true)

	k := api.NewKaiClient()

	require.NotEmpty(t, k.ProcessRegistry())
}
