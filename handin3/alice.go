package handin3


type Alice struct {
	uvw []UVW
}

func (a *Alice) Init(uvw []UVW) {
	a.UVW = uvw
}