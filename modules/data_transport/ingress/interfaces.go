// modules\data_transport\ingress\interfaces.go
package ingress

import "context"

type Handler interface {
	Name() string
	Handle(ctx context.Context, payload []byte) error
}
