package getcreds

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/webdestroya/aws-sso/internal/appconfig"
	"github.com/webdestroya/aws-sso/internal/utils"
	"github.com/webdestroya/aws-sso/internal/utils/awsutils"
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

func writeCliCache(cfg config.SharedConfig, session *config.SSOSession, creds *aws.Credentials) {
	// TODO: write the cache file so the cli can use it
}
