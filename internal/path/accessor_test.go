package path

import "manala/internal/serrors"

func (s *Suite) TestAccessorGetErrors() {
	data := map[string]any{
		"foo": "bar",
		"bar": map[string]any{
			"baz": 123,
		},
	}

	tests := []struct {
		test     string
		path     string
		expected *serrors.Assert
	}{
		{
			test: "Root",
			path: "baz",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "unable to access path",
				Arguments: []any{
					"path", "baz",
				},
			},
		},
		{
			test: "Leaf",
			path: "bar.bar",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "unable to access path",
				Arguments: []any{
					"path", "bar.bar",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			accessor := NewAccessor(
				Path(test.path),
				data,
			)

			value, err := accessor.Get()

			serrors.Equal(s.Assert(), test.expected, err)
			s.Nil(value)
		})
	}
}

func (s *Suite) TestAccessorGet() {
	data := map[string]any{
		"foo": "bar",
		"bar": map[string]any{
			"baz": 123,
		},
	}

	tests := []struct {
		test     string
		path     string
		expected any
	}{
		{
			test:     "Root",
			path:     "foo",
			expected: "bar",
		},
		{
			test:     "Leaf",
			path:     "bar.baz",
			expected: 123,
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			accessor := NewAccessor(
				Path(test.path),
				data,
			)

			value, err := accessor.Get()

			s.NoError(err)
			s.Equal(test.expected, value)
		})
	}
}

func (s *Suite) TestAccessorSet() {
	data := func() map[string]any {
		return map[string]any{
			"foo": "",
			"bar": map[string]any{
				"baz": 0,
			},
		}
	}

	tests := []struct {
		test     string
		path     string
		value    any
		expected map[string]any
	}{
		{
			test:  "Root",
			path:  "foo",
			value: 123,
			expected: map[string]any{
				"foo": 123,
				"bar": map[string]any{
					"baz": 0,
				},
			},
		},
		{
			test:  "Leaf",
			path:  "bar.baz",
			value: "bar",
			expected: map[string]any{
				"foo": "",
				"bar": map[string]any{
					"baz": "bar",
				},
			},
		},
		// See: https://github.com/ohler55/ojg/issues/146
		{
			test:  "Nil",
			path:  "foo",
			value: nil,
			expected: map[string]any{
				"foo": nil,
				"bar": map[string]any{
					"baz": 0,
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			_data := data()

			accessor := NewAccessor(
				Path(test.path),
				_data,
			)

			err := accessor.Set(test.value)

			s.NoError(err)
			s.Equal(test.expected, _data)
		})
	}
}
