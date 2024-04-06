package orderserv

type PacketParams struct{}

func (p PacketParams) GetPrice() int {
	return 5
}

func (p PacketParams) Validate(weight int) bool {
	return weight <= 10
}

type BoxParams struct{}

func (b BoxParams) GetPrice() int {
	return 30
}

func (b BoxParams) Validate(weight int) bool {
	return weight <= 30
}

type TapeParams struct{}

func (t TapeParams) GetPrice() int {
	return 1
}

func (t TapeParams) Validate(weight int) bool {
	return true
}
