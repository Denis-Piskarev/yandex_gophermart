package endpoints

import (
	"github.com/DenisquaP/yandex_gophermart/internal"
)

type Endpoints struct {
	services *internal.ServiceInterface
}

func NewEndpoints(services *internal.ServiceInterface) *Endpoints {
	return &Endpoints{
		services: services,
	}
}
