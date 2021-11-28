package domain

type Result string

// Result
const (
	Success Result = "success"
	Error   Result = "error"
)

type Action string

// Actions (sides)
const (
	Buy  Action = "buy"
	Sell Action = "sell"
	None Action = "none"
)

// Order Types
type OrderType string

const (
	LmtType OrderType = "lmt"
	IocType OrderType = "ioc"
	MktType OrderType = "mkt"
)
