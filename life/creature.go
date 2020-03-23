package life

import "go.stevenxie.me/simulacrum/world"

// A Creature is a living thing.
type Creature interface {
	world.Entity

	// Fuck induces two Creatures to bone.
	//
	// The resulting Creature, offspring, may be nil.
	Fuck(other Creature) (offspring Creature)
}
