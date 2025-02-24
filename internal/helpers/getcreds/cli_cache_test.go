package getcreds_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/aws-sso/internal/helpers/getcreds"
)

func TestCLICacheKey(t *testing.T) {
	tables := []struct {
		ssoName   string
		startUrl  string
		accountId string
		roleName  string
		hash      string
	}{
		{"test-session", "https://blah", "1234567890", "test-role", "f2c4be47e7c208fd9fdb5286a8f46907a665314b"},
		{"", "https://d-92671207e4.awsapps.com/start", "1234567890", "test-role", "048db75bbe50955c16af7aba6ff9c41a3131bb7e"},
	}

	for _, table := range tables {
		t.Run("keytest", func(t *testing.T) {

			sess := &config.SSOSession{
				Name:        table.ssoName,
				SSORegion:   "us-fake-1",
				SSOStartURL: table.startUrl,
			}

			cfg := config.SharedConfig{
				SSORoleName:  table.roleName,
				SSOAccountID: table.accountId,
			}

			// normal
			key := getcreds.CLICacheKey(cfg, sess)
			require.Equal(t, table.hash, key)

			// if the session is blank, but newer style
			cfg2 := config.SharedConfig{
				SSORoleName:  table.roleName,
				SSOAccountID: table.accountId,
				SSOSession:   sess,
			}

			key = getcreds.CLICacheKey(cfg2, nil)
			require.Equal(t, table.hash, key)

			if table.ssoName == "" {
				cfg2 := config.SharedConfig{
					SSORoleName:  table.roleName,
					SSOAccountID: table.accountId,
					SSOStartURL:  table.startUrl,
				}

				key = getcreds.CLICacheKey(cfg2, nil)
				require.Equal(t, table.hash, key)
			} else {
				cfg2 := config.SharedConfig{
					SSORoleName:    table.roleName,
					SSOAccountID:   table.accountId,
					SSOSessionName: table.ssoName,
				}

				key = getcreds.CLICacheKey(cfg2, nil)
				require.Equal(t, table.hash, key)
			}
		})
	}
}
