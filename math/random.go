package math

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

// InitRandom seeds math.rand with crypto/rand (imported as cryptorand), such that future math.rand operations are more or less cryptographically
// secure. It falls back to seeding with current nanosecond time. Without either, the math/rand package will always
// initialize with the same seed (0, I think).
// See: https://stackoverflow.com/a/54491783/5061881
// Imports:
// cryptorand "crypto/rand"
// log "github.com/sirupsen/logrus"
func InitRandom() {
	// Gets 8 bytes using the cryptographically secure random package, and casts them into a uint64 and then an int64
	// (if you use a random byte for the most significant byte of a signed int64 you aren't randomly assigning the sign
	// bit, thus the conversion to unsigned first). I believe it shouldn't matter whether you use LittleEndian or
	// BigEndian, but you need to use one or the other to get to the Uint64([]byte) method.
	var b [8]byte
	_, err := cryptorand.Read(b[:])
	if err != nil {
		log.Warnln("Cannot seed math/rand package with cryptographically secure RNG, using time seed.")
		rand.Seed(time.Now().UTC().UnixNano())
		return
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}
