package awsutils

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
	"github.com/webdestroya/aws-sso/internal/utils"
)

var (
	ErrNotSSOError               = errors.New("no SSO configuration found for profile")
	ErrInvalidSessionConfigError = errors.New("invalid session config")
)

func ExtractSSOInfo(cfg config.SharedConfig) (*config.SSOSession, error) {

	if cfg.SSOSession != nil {
		return cfg.SSOSession, nil
	}

	if len(cfg.SSOStartURL) != 0 && len(cfg.SSORoleName) != 0 && len(cfg.SSOAccountID) != 0 {
		return &config.SSOSession{
			Name:        "",
			SSORegion:   utils.CoalesceString(cfg.SSORegion, cfg.Region, "us-east-1"),
			SSOStartURL: cfg.SSOStartURL,
		}, nil
	}

	return nil, ErrNotSSOError
}

func GetSSOCachePath(session *config.SSOSession) (string, error) {

	if session == nil {
		return "", ErrInvalidSessionConfigError
	}

	cacheKey := utils.CoalesceString(session.Name, session.SSOStartURL)
	if cacheKey == "" {
		return "", fmt.Errorf("%w: missing name or start url", ErrInvalidSessionConfigError)
	}

	tokenFile, err := ssocreds.StandardCachedTokenFilepath(cacheKey)
	if err != nil {
		return "", err
	}
	return tokenFile, nil
}
