package mgollective

type SecurityProvider interface {
	Sign(*Message)
	Verify(Message) bool
}

type SecurityProviderFactory func(*Mgollective) SecurityProvider

var securityProviderRegistry = map[string]SecurityProviderFactory{}

func RegisterSecurityProvider(name string, factory SecurityProviderFactory) {
	securityProviderRegistry[name] = factory
}
