package handin3

type Alice struct {
	UVW []UVW
}

func (a *Alice) Init(uvw []UVW) {
	a.UVW = uvw
}
