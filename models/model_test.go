package models

type mock struct {
	dir string
}

func (mock *mock) getDir() string {
	return mock.dir
}
