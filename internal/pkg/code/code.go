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
	DNSProviderNotSupport = iota + 20000
	DNSProviderSecretError
	DNSProviderExist
)
