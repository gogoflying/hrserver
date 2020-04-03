package types

const (
	STATUS_WAITING = 0
	STATUS_SUCCESS = 1
	STATUS_FAILED  = -1
)

type RepPost struct {
	Name string `json:"my_name"`
}

type RspPost struct {
	Name string `json:"your_name"`
}
