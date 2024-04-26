package code

// user
const (
	UserIsBanned = iota + 10000
	UserNotFound
	UserNotLogin
	UsernameAlreadyExists
)

// dns
const (
	ProviderNotSupport = iota + 20000
	ProviderSecretError
	ProviderExist
	ProviderNotFound
)
