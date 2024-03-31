package orderserv

import (
	order2 "HW1/internal/model/order"
	"math"
)

type PacketParams struct{}

func (p PacketParams) GetParams() order2.PackageParams {
	return order2.PackageParams{
		Price:          5,
		Name:           "packet",
		MaxOrderWeight: 10,
	}
}

type BoxParams struct{}

func (b BoxParams) GetParams() order2.PackageParams {
	return order2.PackageParams{
		Price:          20,
		Name:           "box",
		MaxOrderWeight: 30,
	}
}

type TapeParams struct{}

func (t TapeParams) GetParams() order2.PackageParams {
	return order2.PackageParams{
		Price:          1,
		Name:           "tape",
		MaxOrderWeight: math.MaxInt64,
	}
}
