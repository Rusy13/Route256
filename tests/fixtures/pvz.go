package fixtures

import (
	"Homework/internal/storage/repository"
	"Homework/tests/states"
)

type PvzBuilder struct {
	instance *repository.Pvz
}

func Pvz() *PvzBuilder {
	return &PvzBuilder{instance: &repository.Pvz{}}
}

func (b *PvzBuilder) ID(v int64) *PvzBuilder {
	b.instance.ID = v
	return b
}

func (b *PvzBuilder) Email(v string) *PvzBuilder {
	b.instance.Email = v
	return b
}

func (b *PvzBuilder) PvzName(v string) *PvzBuilder {
	b.instance.PvzName = v
	return b
}

func (b *PvzBuilder) Address(v string) *PvzBuilder {
	b.instance.Address = v
	return b
}

func (b *PvzBuilder) P() *repository.Pvz {
	return b.instance
}

func (b *PvzBuilder) V() repository.Pvz {
	return *b.instance
}

func (b *PvzBuilder) Valid() *PvzBuilder {
	return Pvz().ID(states.Pvz1ID).PvzName(states.Pvz1Name).Email(states.Pvz1Email).Address(states.Pvz1Address)
}
