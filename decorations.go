package dst

type Decorations []string

func (d *Decorations) Add(decs ...string) {
	*d = append(*d, decs...)
}

func (d *Decorations) Replace(decs ...string) {
	*d = decs
}

func (d *Decorations) Clear() {
	*d = nil
}

func (d *Decorations) All() []string {
	return *d
}

type SpaceType int

const (
	None      SpaceType = 0
	NewLine   SpaceType = 1
	EmptyLine SpaceType = 2
)
