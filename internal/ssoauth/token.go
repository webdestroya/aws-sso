package ssoauth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
)

type AwsTokenInfo struct {
	AccessToken string   `json:"accessToken,omitempty"`
	ExpiresAt   *rfc3339 `json:"expiresAt,omitempty"`

	RefreshToken string `json:"refreshToken,omitempty"`
	ClientID     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
}

func ReadTokenInfo(sessionName string) (*AwsTokenInfo, error) {

	tokenPath, err := ssocreds.StandardCachedTokenFilepath(sessionName)
	if err != nil {
		return nil, err
	}

	_ = tokenPath

	return nil, nil
}

func (t *AwsTokenInfo) Expired() bool {
	if t.ExpiresAt == nil {
		return true
	}

	return time.Now().After(time.Time(*t.ExpiresAt))
}

// github.com/aws/aws-sdk-go-v2/credentials@v1.17.26/ssocreds/sso_cached_token.go
type rfc3339 time.Time

func parseRFC3339(v string) (rfc3339, error) {
	parsed, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return rfc3339{}, fmt.Errorf("expected RFC3339 timestamp: %w", err)
	}

	return rfc3339(parsed), nil
}

func (r *rfc3339) UnmarshalJSON(bytes []byte) (err error) {
	var value string

	// Use JSON unmarshal to unescape the quoted value making use of JSON's
	// unquoting rules.
	if err = json.Unmarshal(bytes, &value); err != nil {
		return err
	}

	*r, err = parseRFC3339(value)

	return nil
}

func (r *rfc3339) MarshalJSON() ([]byte, error) {
	value := time.Time(*r).Format(time.RFC3339)

	// Use JSON unmarshal to unescape the quoted value making use of JSON's
	// quoting rules.
	return json.Marshal(value)
}
