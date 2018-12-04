package dst

// TODO: What should we be doing here?
// After cloning a node, should it still be attached to the same object / scope? Or a cloned copy?
// Should we really have objects / scopes at all? As soon as you modify the tree, they are
// potentially invalid.

func CloneObject(o *Object) *Object {
	return nil
}

func CloneScope(s *Scope) *Scope {
	return nil
}
