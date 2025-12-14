package geography

import "math/rand"

// rng is the package-wide pseudo-random number generator used for all
// procedural generation in this package. It can be seeded from main via
// SetRandSeed to ensure reproducible dungeons.
var rng = rand.New(rand.NewSource(1))

// SetRandSeed replaces the package RNG with a new one seeded with the
// provided value. Call this once from main before generating any
// geography to get deterministic output for a given seed.
func SetRandSeed(seed int64) {
	rng = rand.New(rand.NewSource(seed))
}
