package loginflow

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	ssooidcTypes "github.com/aws/aws-sdk-go-v2/service/ssooidc/types"
	"github.com/toqueteos/webbrowser"
	"gopkg.in/ini.v1"

	"github.com/webdestroya/aws-sso/internal/appconfig"
	"github.com/webdestroya/aws-sso/internal/utils"
	"github.com/webdestroya/aws-sso/internal/utils/awsutils"
	"github.com/webdestroya/aws-sso/internal/utils/cmdutils"
)

const (
	pkceGrantType       = `authorization_code`
	pkceChallengeMethod = `S256`
	defaultSSOScope     = `sso:account:access`
	pkceCallbackTimeout = 180 * time.Second
	pkceCallbackPath    = `/oauth/callback`
	pkceShutdownTimeout = 5 * time.Second
	ssoRegScopesKey     = `sso_registration_scopes`
)

type callbackResult struct {
	code    string
	state   string
	authErr string
}

//go:embed success.html
var pkceSuccessHTML []byte

//go:embed error.html
var pkceErrorHTML []byte

func doLoginFlowPKCE(ctx context.Context, out io.Writer, session *config.SSOSession, options *loginFlowOptions) (*awsutils.AwsTokenInfo, error) {

	tokenFile, err := awsutils.GetSSOCachePath(session)
	if err != nil {
		return nil, err
	}

	cfg, err := awsutils.LoadDefaultConfig(ctx, config.WithRegion(session.SSORegion))
	if err != nil {
		return nil, err
	}

	client := cmdutils.NewSSOOIDCClient(cfg)

	verifier, err := generateCodeVerifier()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PKCE code verifier: %w", err)
	}
	challenge := computeCodeChallenge(verifier)

	state, err := generateState()
	if err != nil {
		return nil, fmt.Errorf("failed to generate OAuth state: %w", err)
	}

	scopes := resolveRegistrationScopes(session)

	ln, resultCh, srv, err := startCallbackServer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start local callback server: %w", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), pkceShutdownTimeout)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	go func() {
		// http.ErrServerClosed is expected during shutdown
		if serveErr := srv.Serve(ln); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			// Nothing actionable here; the main flow will time out or observe ctx cancellation.
			_ = serveErr
		}
	}()

	port := ln.Addr().(*net.TCPAddr).Port
	redirectUri := fmt.Sprintf("http://127.0.0.1:%d%s", port, pkceCallbackPath)

	regResp, err := client.RegisterClient(ctx, &ssooidc.RegisterClientInput{
		ClientName:   aws.String(clientName),
		ClientType:   aws.String(clientType),
		GrantTypes:   []string{`authorization_code`, `refresh_token`},
		IssuerUrl:    aws.String(session.SSOStartURL),
		RedirectUris: []string{redirectUri},
		Scopes:       scopes,
	})
	if err != nil {
		return nil, err
	}

	authEndpoint := fmt.Sprintf(`https://oidc.%s.amazonaws.com/authorize`, session.SSORegion)

	authURL, err := url.Parse(authEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse authorization endpoint: %w", err)
	}

	q := authURL.Query()
	q.Set("response_type", "code")
	q.Set("client_id", aws.ToString(regResp.ClientId))
	q.Set("redirect_uri", redirectUri)
	q.Set("code_challenge", challenge)
	q.Set("code_challenge_method", pkceChallengeMethod)
	q.Set("scopes", strings.Join(scopes, " "))
	q.Set("state", state)
	authURL.RawQuery = q.Encode()

	authURLStr := authURL.String()

	fmt.Fprintln(out)
	fmt.Fprintf(out, "Logging in to: %s\n", utils.CoalesceString(session.Name, session.SSOStartURL))
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Attempting to automatically open the SSO authorization page in your default browser.")
	fmt.Fprintln(out, "If the browser does not open, open the following URL manually:")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "  "+authURLStr)
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Waiting for authentication...")

	if !options.DisableBrowser {
		go func() {
			if berr := webbrowser.Open(authURLStr); berr != nil {
				fmt.Fprintln(out)
				fmt.Fprintln(out, utils.WarningStyle.Render(fmt.Sprintf("WARNING: Failed to open browser URL: %v", berr.Error())))
				fmt.Fprintln(out, "You will need to visit the authorization URL manually.")
				fmt.Fprintln(out)
			}
		}()
	}

	var code string
	select {
	case res := <-resultCh:
		if res.authErr != "" {
			return nil, cmdutils.NewNonUsageError(fmt.Sprintf("authorization failed: %s", res.authErr))
		}
		if res.state != state {
			return nil, ErrAuthDeniedError
		}
		if res.code == "" {
			return nil, ErrAuthDeniedError
		}
		code = res.code
	case <-time.After(pkceCallbackTimeout):
		return nil, ErrAuthTimeoutError
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	createTokenOut, err := client.CreateToken(ctx, &ssooidc.CreateTokenInput{
		ClientId:     regResp.ClientId,
		ClientSecret: regResp.ClientSecret,
		GrantType:    aws.String(pkceGrantType),
		Code:         aws.String(code),
		CodeVerifier: aws.String(verifier),
		RedirectUri:  aws.String(redirectUri),
	})
	if err != nil {
		if _, ok := errors.AsType[*ssooidcTypes.ExpiredTokenException](err); ok {
			return nil, ErrAuthTimeoutError
		}

		if _, ok := errors.AsType[*ssooidcTypes.AccessDeniedException](err); ok {
			return nil, ErrAuthDeniedError
		}

		return nil, err
	}

	expiresAt := awsutils.RFC3339(time.Now().Add(time.Duration(createTokenOut.ExpiresIn) * time.Second))

	tokenObj := &awsutils.AwsTokenInfo{
		AccessToken:  aws.ToString(createTokenOut.AccessToken),
		ExpiresAt:    &expiresAt,
		RefreshToken: aws.ToString(createTokenOut.RefreshToken),
		ClientID:     aws.ToString(regResp.ClientId),
		ClientSecret: aws.ToString(regResp.ClientSecret),
	}

	fmt.Fprintln(out, "Authentication successful! Writing token...")
	fmt.Fprintln(out)

	if err := awsutils.WriteAWSToken(tokenFile, *tokenObj); err != nil {
		return nil, err
	}

	return tokenObj, nil
}

// generateCodeVerifier returns a high-entropy, base64url (no padding) PKCE
// code verifier, 43 characters long (derived from 32 random bytes).
func generateCodeVerifier() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// computeCodeChallenge returns base64url(sha256(verifier)) with no padding,
// per RFC 7636 for the S256 method.
func computeCodeChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

// generateState returns a random hex-encoded OAuth state value used for
// CSRF protection on the redirect callback.
func generateState() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// resolveRegistrationScopes looks up sso_registration_scopes from the
// [sso-session <name>] section of ~/.aws/config. Falls back to a single
// scope of defaultSSOScope when the value is unset, the session has no
// name, or the config file cannot be read.
func resolveRegistrationScopes(session *config.SSOSession) []string {
	fallback := []string{defaultSSOScope}
	if session == nil || session.Name == "" {
		return fallback
	}

	f, err := utils.LoadIniFile(ini.LoadOptions{
		SkipUnrecognizableLines: true,
		Insensitive:             true,
		Loose:                   true,
	}, appconfig.GetAwsConfigPath())
	if err != nil || f == nil {
		return fallback
	}

	sec, err := f.GetSection("sso-session " + session.Name)
	if err != nil {
		return fallback
	}

	key, err := sec.GetKey(ssoRegScopesKey)
	if err != nil {
		return fallback
	}

	raw := strings.TrimSpace(key.String())
	if raw == "" {
		return fallback
	}

	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}
	if len(out) == 0 {
		return fallback
	}
	return out
}

// startCallbackServer binds an HTTP listener to 127.0.0.1 on an
// OS-assigned port and returns the listener, a buffered result channel
// (cap 1), and the server. The caller is responsible for invoking
// srv.Serve and srv.Shutdown.
func startCallbackServer(ctx context.Context) (net.Listener, chan callbackResult, *http.Server, error) {
	lc := &net.ListenConfig{}
	ln, err := lc.Listen(ctx, "tcp", "127.0.0.1:0")
	if err != nil {
		return nil, nil, nil, err
	}

	resultCh := make(chan callbackResult, 1)

	mux := http.NewServeMux()
	mux.HandleFunc(pkceCallbackPath, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		res := callbackResult{
			code:    q.Get("code"),
			state:   q.Get("state"),
			authErr: q.Get("error"),
		}
		if desc := q.Get("error_description"); desc != "" && res.authErr != "" {
			res.authErr = res.authErr + ": " + desc
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if res.authErr != "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, string(pkceErrorHTML), res.authErr)
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(pkceSuccessHTML)
		}

		select {
		case resultCh <- res:
		default:
		}
	})

	srv := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return ln, resultCh, srv, nil
}
