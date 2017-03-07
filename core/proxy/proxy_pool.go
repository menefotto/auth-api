package proxy

import "math/rand"

type Pool struct {
	Connections map[*UsersVerifier]*int
}

func NewPool(size int) *Pool {
	p := &Pool{make(map[*UsersVerifier]*int, size)}

	for i := 0; i < size; i++ {
		u := New()
		n := rand.Int()
		p.Connections[u] = &n
	}

	return p
}

func (p *Pool) Get() *UsersVerifier {

	for k, v := range p.Connections {
		if *v != 0 {
			*v = rand.Int()

			return k
		}
	}

	return nil
}

func (p *Pool) Put(u *UsersVerifier) {

	for _, v := range p.Connections {
		if *v == 0 {
			*v = rand.Int()

			return
		}
	}
}
