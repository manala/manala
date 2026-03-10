package path

import (
	"github.com/manala/manala/internal/serrors"

	"github.com/ohler55/ojg/jp"
)

type Accessor struct {
	path Path
	data any
}

func NewAccessor(path Path, data any) Accessor {
	return Accessor{
		path: path,
		data: data,
	}
}

func (accessor Accessor) Get() (any, error) {
	expr, err := accessor.expr()
	if err != nil {
		return nil, err
	}

	value, found := expr.FirstFound(accessor.data)
	if !found {
		return nil, serrors.New("unable to access path").
			WithArguments("path", accessor.path.String())
	}

	return value, nil
}

func (accessor Accessor) Set(value any) error {
	expr, err := accessor.expr()
	if err != nil {
		return err
	}

	return expr.Set(accessor.data, value)
}

func (accessor Accessor) expr() (jp.Expr, error) {
	return jp.ParseString(accessor.path.String())
}
