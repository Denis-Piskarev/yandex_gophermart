package internal

// Service - struct of service layer interfaces
type Service struct {
	BalanceKeeper
	AuthInterface
	OrderInterface
	TokenInterface
}
