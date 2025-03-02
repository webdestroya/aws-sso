package awsutils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
	"github.com/webdestroya/aws-sso/internal/utils"
)

type AwsTokenInfo struct {
	AccessToken string   `json:"accessToken,omitempty"`
	ExpiresAt   *RFC3339 `json:"expiresAt,omitempty"`

	RefreshToken string `json:"refreshToken,omitempty"`
	ClientID     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
}

// SESSION NAME IS NOT PROFILE NAME
func ReadTokenInfo(sessionName string) (*AwsTokenInfo, error) {
	tokenPath, err := ssocreds.StandardCachedTokenFilepath(sessionName)
	if err != nil {
		return nil, err
	}

	return ReadTokenFile(tokenPath)
}

func ReadTokenFile(tokenFilePath string) (*AwsTokenInfo, error) {
	data, err := os.ReadFile(tokenFilePath)
	if err != nil {
		return nil, err
	}

	token := &AwsTokenInfo{}

	if err := json.Unmarshal(data, token); err != nil {
		return nil, err
	}

	return token, nil
}

func (t *AwsTokenInfo) Expired() bool {
	if t.ExpiresAt == nil {
		return true
	}

	return time.Now().After(time.Time(*t.ExpiresAt))
}

// github.com/aws/aws-sdk-go-v2/credentials@v1.17.26/ssocreds/sso_cached_token.go
type RFC3339 time.Time

func (r RFC3339) String() string {
	return r.AsTime().String()
}

func (r RFC3339) AsTime() time.Time {
	return time.Time(r)
}

func parseRFC3339(v string) (RFC3339, error) {
	parsed, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return RFC3339{}, fmt.Errorf("expected RFC3339 timestamp: %w", err)
	}

	return RFC3339(parsed), nil
}

func (r *RFC3339) UnmarshalJSON(bytes []byte) (err error) {
	var value string

	// Use JSON unmarshal to unescape the quoted value making use of JSON's
	// unquoting rules.
	if err = json.Unmarshal(bytes, &value); err != nil {
		return err
	}

	*r, err = parseRFC3339(value)

	return nil
}

func (r *RFC3339) MarshalJSON() ([]byte, error) {
	value := time.Time(*r).Format(time.RFC3339)

	// Use JSON unmarshal to unescape the quoted value making use of JSON's
	// quoting rules.
	return json.Marshal(value)
}

func WriteAWSToken(filename string, t AwsTokenInfo) error {
	return storeCachedToken(filename, t, 0600)
}

func storeCachedToken(filename string, t AwsTokenInfo, fileMode os.FileMode) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return utils.AtomicWriteFile(filename, data, fileMode)
}

/*
func storeCachedToken(filename string, t AwsTokenInfo, fileMode os.FileMode) (err error) {
	// tmpFilename := filename + ".tmp-" + strconv.FormatInt(sdk.NowTime().UnixNano(), 10)
	tmpFilename := filename + ".tmp-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	if err := writeCacheFile(tmpFilename, fileMode, t); err != nil {
		return err
	}

	if err := os.Rename(tmpFilename, filename); err != nil {
		return fmt.Errorf("failed to replace old cached SSO token file, %w", err)
	}

	return nil
}

func writeCacheFile(filename string, fileMode os.FileMode, t AwsTokenInfo) (err error) {
	var f *os.File
	f, err = os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, fileMode)
	if err != nil {
		return fmt.Errorf("failed to create cached SSO token file %w", err)
	}

	defer func() {
		closeErr := f.Close()
		if err == nil && closeErr != nil {
			err = fmt.Errorf("failed to close cached SSO token file, %w", closeErr)
		}
	}()

	encoder := json.NewEncoder(f)

	if err = encoder.Encode(t); err != nil {
		return fmt.Errorf("failed to serialize cached SSO token, %w", err)
	}

	return nil
}
*/
