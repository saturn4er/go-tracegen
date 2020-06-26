package example

import (
	"context"

	"github.com/saturn4er/go-traceutil"
	"go.uber.org/multierr"
)

var _ = multierr.Append
var _ = traceutil.ChildSpanFromContext

type _ = context.Context

type SomeStructTracer struct {
	next *SomeStruct
}

func (t *SomeStructTracer) WithContext(ctx context.Context, param1 string, params2 int) {
	ctx, span := traceutil.ChildSpanFromContext(ctx, "SomeStruct.Get")

	defer span.Finish()

	t.next.WithContext(ctx, param1, params2)
}
func (t *SomeStructTracer) WithContextAndReturn(ctx context.Context, param1 string, params2 int) (_ int) {
	ctx, span := traceutil.ChildSpanFromContext(ctx, "SomeStruct.Get")

	defer span.Finish()

	return t.next.WithContextAndReturn(ctx, param1, params2)
}
func (t *SomeStructTracer) WithContextAndReturnAndError(ctx context.Context, param1 string, params2 int) (_ int, rerr1 error) {
	ctx, span := traceutil.ChildSpanFromContext(ctx, "SomeStruct.Get")

	defer func() { traceutil.FinishSpanWithErr(span, multierr.Combine(rerr1)) }()

	return t.next.WithContextAndReturnAndError(ctx, param1, params2)
}
func (t *SomeStructTracer) WithContextAndError(ctx context.Context, param1 string, params2 int) (rerr0 error) {
	ctx, span := traceutil.ChildSpanFromContext(ctx, "SomeStruct.Get")

	defer func() { traceutil.FinishSpanWithErr(span, multierr.Combine(rerr0)) }()

	return t.next.WithContextAndError(ctx, param1, params2)
}
func (t *SomeStructTracer) WithEllipsis(ctx context.Context, param1 ...int) {
	ctx, span := traceutil.ChildSpanFromContext(ctx, "SomeStruct.Get")

	defer span.Finish()

	t.next.WithEllipsis(ctx, param1...)
}

func NewSomeStructTracer(next *SomeStruct) *SomeStructTracer {
	return &SomeStructTracer{
		next: next,
	}
}

type SomeInterfaceTracer struct {
	next SomeInterface
}

func (t *SomeInterfaceTracer) WithContext(ctx context.Context, param1 string, params2 int) {
	ctx, span := traceutil.ChildSpanFromContext(ctx, "SomeInterface.Get")

	defer span.Finish()

	t.next.WithContext(ctx, param1, params2)
}
func (t *SomeInterfaceTracer) WithContextAndReturn(ctx context.Context, param1 string, params2 int) (_ int) {
	ctx, span := traceutil.ChildSpanFromContext(ctx, "SomeInterface.Get")

	defer span.Finish()

	return t.next.WithContextAndReturn(ctx, param1, params2)
}
func (t *SomeInterfaceTracer) WithContextAndReturnAndError(ctx context.Context, param1 string, params2 int) (_ int, rerr1 error) {
	ctx, span := traceutil.ChildSpanFromContext(ctx, "SomeInterface.Get")

	defer func() { traceutil.FinishSpanWithErr(span, multierr.Combine(rerr1)) }()

	return t.next.WithContextAndReturnAndError(ctx, param1, params2)
}

func (t *SomeInterfaceTracer) WithContextAndError(ctx context.Context, param1 string, params2 int) (rerr0 error) {
	ctx, span := traceutil.ChildSpanFromContext(ctx, "SomeInterface.Get")

	defer func() { traceutil.FinishSpanWithErr(span, multierr.Combine(rerr0)) }()

	return t.next.WithContextAndError(ctx, param1, params2)
}

func (t *SomeInterfaceTracer) WithEllipsis(ctx context.Context, param1 ...int) {
	ctx, span := traceutil.ChildSpanFromContext(ctx, "SomeInterface.Get")

	defer span.Finish()

	t.next.WithEllipsis(ctx, param1...)
}

func NewSomeInterfaceTracer(next SomeInterface) *SomeInterfaceTracer {
	return &SomeInterfaceTracer{
		next: next,
	}
}
