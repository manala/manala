package serror

func New(msg string) Error {
	return Error{
		msg: msg,
	}
}

type Error struct {
	msg   string
	attrs [][2]any
	dump  string
	err   error
}

func (err Error) Error() string {
	return err.msg
}

func (err Error) Attrs() [][2]any {
	return err.attrs
}

func (err Error) With(args ...any) Error {
	a := args
	for len(a) >= 2 {
		if key, ok := a[0].(string); ok {
			err.attrs = append(err.attrs, [2]any{key, a[1]})
		}
		a = a[2:]
	}
	return err
}

func (err Error) Dump() string { return err.dump }

func (err Error) WithDump(dump string) Error {
	err.dump = dump
	return err
}

func (err Error) Err() error {
	return err.err
}

func (err Error) WithErr(e error) Error {
	err.err = e
	return err
}

type Attrs interface {
	Attrs() [][2]any
}
