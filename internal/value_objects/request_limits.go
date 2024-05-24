package value_objects

type RequestLimits struct {
	IPLimit       uint
	APILimit      uint
	ExpireSeconds uint
}

func NewRequestLimit(IPLimit uint, APILimit uint, ExpireSeconds uint) RequestLimits {
	return RequestLimits{
		IPLimit:       IPLimit,
		APILimit:      APILimit,
		ExpireSeconds: ExpireSeconds,
	}
}
