package di

import (
	"errors"
)

// Builder is a Container builder
type Builder struct {
	ctn *Container
}

// Add adds Def to the Container
func (b *Builder) Add(defs ...Def) error {
	b.initContainer()
	for _, def := range defs {
		if _, exists := b.ctn.defs[def.Name]; exists {
			return errors.New("[bubucore.di] definition with name `" + def.Name + "` already exists")
		}
		if def.Validate != nil {
			if err := def.Validate(b.ctn); err != nil {
				return err
			}
		}
		b.ctn.defs[def.Name] = def
	}
	return nil
}

// Build prepares Container and builds non-lazy definitions
func (b *Builder) Build() (*Container, error) {
	b.initContainer()
	for name, def := range b.ctn.defs {
		if !def.Lazy {
			if err := def.build(b.ctn); err != nil {
				return nil, err
			}
			b.ctn.defs[name] = def
		}
	}
	return b.ctn, nil
}

// initContainer creates container instance
func (b *Builder) initContainer() {
	if b.ctn == nil {
		b.ctn = &Container{
			defs: make(map[string]Def),
		}
	}
}
