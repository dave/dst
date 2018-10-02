package dst

type Decorations []string

func (d *Decorations) Add(text string) {
	*d = append(*d, text)
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
