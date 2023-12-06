package common

const (
	FirstVersion = 1
)

const (
	TransactionStatusPending = 1 << iota
	TransactionStatusRejected
	TransactionStatusSuccess
)

const (
	WalletStatusActive = 1 << iota
)

const (
	UserStatusActive = 1 << iota
)
const (
	GoldenDollar = "GDD"
)

const (
	TopupStatusSuccess = 1 << iota
)
