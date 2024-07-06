package ctx

import (
	"context"
	"fmt"
)

type Ctxer struct{}

func (c Ctxer) Ctx(ctx context.Context, context, namespace string) error {
	fmt.Println("ctx called", context, namespace)
	return nil
}
