package interests

type Query struct {
	Limit   uint32
	Sort    Sort
	Order   Order
	Pattern string
	Public  bool // include public non-own?
}
