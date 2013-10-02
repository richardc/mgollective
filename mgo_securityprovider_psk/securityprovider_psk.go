package mgo_securityprovider_psk

import (
	"crypto/md5"
	"fmt"
	"github.com/richardc/mgollective/mgollective"
	"io"
)

type PskSecurityProvider struct {
	psk string
}

// The psk provider md5 hashes a concatenation of the psk and the message body
// and returns the lowercased hex string for that

func (s PskSecurityProvider) hash(message []byte) string {
	hash := md5.New()
	io.WriteString(hash, s.psk)
	hash.Write(message)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (p PskSecurityProvider) Sign(message []byte) map[string]string {
	headers := make(map[string]string)
	headers["mc_signature"] = p.hash(message)
	return headers
}

func (p PskSecurityProvider) Verify(message []byte, headers map[string]string) bool {
	if headers["mc_signature"] == p.hash(message) {
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
