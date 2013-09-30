package mgollective

type SecurityProvider interface {
	Sign([]byte) map[string]string
	Verify([]byte, map[string]string) bool
}

type SecurityProviderFactory func(*Mgollective) SecurityProvider

var securityProviderRegistry = map[string]SecurityProviderFactory{}

func RegisterSecurityProvider(name string, factory SecurityProviderFactory) {
	securityProviderRegistry[name] = factory
}
