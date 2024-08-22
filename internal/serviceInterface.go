package internal

// Service - str of service layer interfaces
type Service struct {
	BalanceKeeper
	AuthInterface
	OrderInterface
	TokenInterface
}
