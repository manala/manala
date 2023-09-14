package path

type Path string

func (p Path) String() string {
	return string(p)
}

func (p Path) Join(seg string) Path {
	if p != "" {
		seg = "." + seg
	}
	return Path(string(p) + seg)
}
