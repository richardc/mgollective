package mgo_securityprovider_psk

import (
	"github.com/richardc/mgollective/mgollective"
)

type PskSecurityProvider struct {
	psk string
}

func (p PskSecurityProvider) Sign(message *mgollective.Message) {
}

func (p PskSecurityProvider) Verify(message mgollective.Message) bool {
	return true
}

func makePskSecurityProvider(app *mgollective.Mgollective) mgollective.SecurityProvider {
	return &PskSecurityProvider{
		psk: "pies",
	}
}

func init() {
	mgollective.DeclareConfig("plugin.psk")
	mgollective.RegisterSecurityProvider("psk", makePskSecurityProvider)
}
