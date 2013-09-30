package mgo_securityprovider_psk

import (
	"github.com/richardc/mgollective/mgollective"
)

type PskSecurityProvider struct {
	psk string
}

func (p PskSecurityProvider) Sign(message []byte) map[string]string {
	headers := make(map[string]string)
	headers["signature"] = p.psk
	return headers
}

func (p PskSecurityProvider) Verify(message []byte, headers map[string]string) bool {
	if headers["signature"] == p.psk {
		return true
	}
	return false
}

func makePskSecurityProvider(app *mgollective.Mgollective) mgollective.SecurityProvider {
	return &PskSecurityProvider{
		psk: app.GetConfig("plugin.psk", ""),
	}
}

func init() {
	mgollective.DeclareConfig("plugin.psk")
	mgollective.RegisterSecurityProvider("psk", makePskSecurityProvider)
}
