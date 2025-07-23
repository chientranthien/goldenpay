package common

const (
	FirstVersion = 1
)

const (
	TransactionStatusPending  = 1
	TransactionStatusRejected = 2
	TransactionStatusSuccess  = 3
)

const (
	WalletStatusActive = 1
)

const (
	UserStatusActive = 1
)

const (
	ContactStatusActive = 1
)

const (
	GoldenDollar = "GDD"
)

const (
	TopupStatusSuccess = 1
)

const (
	// Channel statuses
	ChannelStatusActive   = 1
	ChannelStatusArchived = 2
	ChannelStatusDeleted  = 3

	// Message statuses
	MessageStatusActive  = 1
	MessageStatusEdited  = 2
	MessageStatusDeleted = 3
)

const (
	// Channel types
	ChannelTypePublic  = "public"
	ChannelTypePrivate = "private"
	ChannelTypeDirect  = "direct"

	// Message types
	MessageTypeText   = "text"
	MessageTypeFile   = "file"
	MessageTypeSystem = "system"

	// User presence statuses
	PresenceStatusOnline  = "online"
	PresenceStatusAway    = "away"
	PresenceStatusBusy    = "busy"
	PresenceStatusOffline = "offline"
)
