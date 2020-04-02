package gen_mmo

import (
	"github.com/itfantasy/gonode"
)

var (
	Err_PeerNotFound              error = gonode.CustomError(1, "PeerNotFound")
	Err_InvalidOperation                = gonode.CustomError(54, "Err_InvalidOperation")
	Err_ItemAccessDenied                = gonode.CustomError(55, "ItemAccessDenied")
	Err_InterestAreaNotFound            = gonode.CustomError(56, "InterestAreaNotFound")
	Err_InterestAreaAlreadyExists       = gonode.CustomError(57, "InterestAreaAlreadyExists")
	Err_WorldAlreadyExists              = gonode.CustomError(101, "WorldAlreadyExists")
	Err_WorldNotFound                   = gonode.CustomError(102, "WorldNotFound")
	Err_ItemAlreadyExists               = gonode.CustomError(103, "ItemAlreadyExists")
	Err_ItemNotFound                    = gonode.CustomError(104, "ItemNotFound")
)
