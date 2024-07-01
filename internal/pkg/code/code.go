package code

// user
const (
	UserIsBanned = iota + 10000
	UserNotFound
	UserNotLogin
	UsernameAlreadyExists
)

// provider
const (
	ProviderNotSupport = iota + 20000
	ProviderSecretError
	ProviderExist
	ProviderNotFound
)

// domain
const (
	DomainNotFound = iota + 30000
	DomainIsManaged
	RecordCreateFailed
	RecordUpdateFailed
)

// cert
const (
	InvalidDomain = iota + 40000
	CertNotFound
	CertNotActive
	CertStatusApply
)

// cdn
const (
	CDNNotFound = iota + 50000
)

// cert hosting
const (
	CertHostingNotFound = iota + 60000
	CertHostingExist
	CertHostingDeploy
	CertHostingTypeError
)
