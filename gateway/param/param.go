package param

import pb "github.com/Far-sa/commons/api"

type CreateOrderResponse struct {
	Order         *pb.Order `json:"order"`
	RedirectToURL string    `json:"redirectToURL"`
}
