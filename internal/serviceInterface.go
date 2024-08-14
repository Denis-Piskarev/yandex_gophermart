package internal

// ServiceInterface - интерфейс слоя service
type ServiceInterface struct {
	BalanceKeeper
	AuthInterface
	OrderInterface
	TokenInterface
}
