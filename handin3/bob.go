package handin3

type Bob struct {
	uvw []UVW
}


func (b *Bob) Init(uvw []UVW) {
	b.uvw = uvw
}