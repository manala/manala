package errors

// At creates an Error positioned at the given offset.
func At(err error, src string, offset int64) Error {
	e := Error{
		error:  err,
		line:   0,
		column: 0,
	}

	if src == "" {
		return e
	}

	// Compute position
	e.line, e.column = 1, 1
	for _, r := range src[:offset-1] {
		if r == '\n' {
			e.line++
			e.column = 1
		} else {
			e.column++
		}
	}

	return e
}
