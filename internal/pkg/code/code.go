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
)
