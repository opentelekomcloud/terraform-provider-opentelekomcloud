package v3

const (
	keyClient       = "lbv3-client"
	ErrCreateClient = "error creating ELBv3 client: %w"
)

func iBool(v interface{}) *bool {
	b := v.(bool)
	return &b
}
