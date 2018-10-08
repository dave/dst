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

type DecorationStmtDecorations struct {
	Start Decorations
}

type DecorationDeclDecorations struct {
	Start Decorations
}

type SpaceType int

const (
	None      SpaceType = 0
	NewLine   SpaceType = 1
	EmptyLine SpaceType = 2
)
