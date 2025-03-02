package getcreds

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
	"github.com/webdestroya/aws-sso/internal/appconfig"
	"github.com/webdestroya/aws-sso/internal/utils"
	"github.com/webdestroya/aws-sso/internal/utils/awsutils"
)

type cliCacheCredentials struct {
	AccessKeyId     string `json:",omitempty"`
	SecretAccessKey string `json:",omitempty"`
	SessionToken    string `json:",omitempty"`
	Expiration      string `json:",omitempty"`
}

type cliCacheEntry struct {
	ProviderType string              `json:",omitempty"`
	Credentials  cliCacheCredentials `json:",omitempty"`
}

var (
	ErrInvalidCacheEntryError = errors.New("invalid cli cache entry")
)

func CLICacheKey(cfg config.SharedConfig, session *config.SSOSession) string {

	if session == nil {
		if exSess, err := awsutils.ExtractSSOInfo(cfg); err == nil {
			session = exSess
		} else {
			session = &config.SSOSession{}
		}
	}

	parts := make([]string, 0, 3)

	parts = append(parts, fmt.Sprintf(`%q:%q`, "accountId", cfg.SSOAccountID))
	parts = append(parts, fmt.Sprintf(`%q:%q`, "roleName", cfg.SSORoleName))

	sessName := utils.CoalesceString(session.Name, cfg.SSOSessionName)

	if sessName != "" {
		parts = append(parts, fmt.Sprintf(`%q:%q`, "sessionName", sessName))
	} else {
		startUrl := utils.CoalesceString(session.SSOStartURL, cfg.SSOStartURL)
		parts = append(parts, fmt.Sprintf(`%q:%q`, "startUrl", startUrl))
	}

	slices.Sort(parts)

	input := "{" + strings.Join(parts, `,`) + "}"

	sum := sha1.Sum([]byte(input))

	return hex.EncodeToString(sum[:])
}

func CLICacheFile(cfg config.SharedConfig, session *config.SSOSession) string {
	return filepath.Join(
		filepath.Dir(appconfig.GetAwsConfigPath()),
		"cli",
		"cache",
		CLICacheKey(cfg, session)+".json",
	)
}

func writeCliCache(cfg config.SharedConfig, session *config.SSOSession, creds *aws.Credentials) error {
	obj := cliCacheEntry{
		ProviderType: "sso",
		Credentials: cliCacheCredentials{
			AccessKeyId:     creds.AccessKeyID,
			SecretAccessKey: creds.SecretAccessKey,
			SessionToken:    creds.SessionToken,
			Expiration:      creds.Expires.UTC().Format(time.RFC3339),
		},
	}

	jsonOut, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	return utils.WriteFile(CLICacheFile(cfg, session), jsonOut, 0600)

}

func ReadCliCache(cfg config.SharedConfig, session *config.SSOSession) (*aws.Credentials, error) {
	obj := cliCacheEntry{}

	cacheFile := CLICacheFile(cfg, session)

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	if obj.ProviderType != "sso" || obj.Credentials.SecretAccessKey == "" || obj.Credentials.AccessKeyId == "" {
		return nil, ErrInvalidCacheEntryError
	}

	expires, err := time.Parse(time.RFC3339, obj.Credentials.Expiration)
	if err != nil {
		return nil, err
	}

	creds := &aws.Credentials{
		AccessKeyID:     obj.Credentials.AccessKeyId,
		SecretAccessKey: obj.Credentials.SecretAccessKey,
		SessionToken:    obj.Credentials.SessionToken,
		Source:          ssocreds.ProviderName,
		CanExpire:       true,
		Expires:         expires,
		AccountID:       cfg.SSOAccountID,
	}
	return creds, nil

}
