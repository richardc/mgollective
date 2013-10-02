package mgo_securityprovider_null

import (
	"github.com/richardc/mgollective/mgollective"
)

type NullSecurityProvider struct {
}

func (p NullSecurityProvider) Sign(message []byte) map[string]string {
	headers := make(map[string]string)
	headers["mc_signature"] = "Everything looks fine buddy"
	return headers
}

func (p NullSecurityProvider) Verify(message []byte, headers map[string]string) bool {
	return true
}

func makeNullSecurityProvider(app *mgollective.Mgollective) mgollective.SecurityProvider {
	return &NullSecurityProvider{}
}

func init() {
	mgollective.RegisterSecurityProvider("null", makeNullSecurityProvider)
}
