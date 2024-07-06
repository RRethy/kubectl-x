package ctx

import (
	"context"
)

func Ctx(ctx context.Context, context, namespace string) error {
	return Ctxer{}.Ctx(ctx, context, namespace)
}
