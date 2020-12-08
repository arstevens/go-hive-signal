package negotiator

type UnmarshalNegotiateMessage func(b []byte) (interface{}, error)
type NegotiateMessage interface {
	IsAccepted() bool
}
