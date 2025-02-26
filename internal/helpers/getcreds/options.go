package getcreds

import "github.com/webdestroya/aws-sso/internal/helpers/loginflow"

type getCredOptions struct {
	DisableLogin     bool
	CliCache         bool
	LoginFlowOptions []loginflow.LoginFlowOption
}

type GetCredOption func(*getCredOptions)

func WithLoginFlowOptions(lfo ...loginflow.LoginFlowOption) GetCredOption {
	return func(o *getCredOptions) {
		o.LoginFlowOptions = lfo
	}
}

func WithLoginDisabled() GetCredOption {
	return func(o *getCredOptions) {
		o.DisableLogin = true
	}
}

func WithCliCache(v bool) GetCredOption {
	return func(o *getCredOptions) {
		o.CliCache = v
	}
}
