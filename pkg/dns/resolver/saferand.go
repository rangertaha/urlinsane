package resolver

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	math_rand "math/rand"
	"sync"
	"time"
)

type saferand struct {
	sync.Mutex

	rand *math_rand.Rand
}

func newSafeRand(crypto bool) *saferand {
	var source math_rand.Source

	if crypto {
		var b [8]byte
		_, err := crypto_rand.Read(b[:])

		if err == nil {
			source = math_rand.NewSource(int64(binary.LittleEndian.Uint64(b[:])))
		}
	}

	if source == nil {
		source = math_rand.NewSource(time.Now().UTC().UnixNano())
	}

	mathrand := math_rand.New(source)
	saferand := saferand{rand: mathrand}

	return &saferand
}

func (s *saferand) Int31() int32 {
	s.Lock()
	defer s.Unlock()

	return s.rand.Int31()
}
