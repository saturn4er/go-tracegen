package example

import (
	"context"
)

type SomeStruct struct {
}

func (s *SomeStruct) WithContext(ctx context.Context, param1 string, params2 int) {

}
func (s *SomeStruct) WithContextAndReturn(ctx context.Context, param1 string, params2 int) int {
	return 1
}
func (s *SomeStruct) WithContextAndReturnAndError(ctx context.Context, param1 string, params2 int) (int, error) {
	return 0, nil
}
func (s *SomeStruct) WithContextAndError(ctx context.Context, param1 string, params2 int) error {
	return nil
}
func (s *SomeStruct) WithEllipsis(ctx context.Context, param1 ...int) {
}

type SomeInterface interface {
	WithContext(ctx context.Context, param1 string, params2 int)
	WithContextAndReturn(ctx context.Context, param1 string, params2 int) int
	WithContextAndReturnAndError(ctx context.Context, param1 string, params2 int) (int, error)
	WithContextAndError(ctx context.Context, param1 string, params2 int) error
	WithEllipsis(ctx context.Context, param1 ...int)
}
