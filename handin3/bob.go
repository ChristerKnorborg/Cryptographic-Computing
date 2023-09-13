package handin3

type Bob struct {
	UVW []UVW
}

func (b *Bob) Init(uvw []UVW) {
	b.UVW = uvw
}
