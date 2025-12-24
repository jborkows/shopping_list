package admin

import "context"

type Maintenance interface {
	OptimizeDB(ctx context.Context) error
}
