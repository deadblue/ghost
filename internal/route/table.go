package route

type _StrictKey struct {
	method, path string
}

type _StrictTable[T any] map[_StrictKey]T

func (t _StrictTable[T]) Put(method, path string, target T) {
	key := _StrictKey{
		method: method,
		path:   path,
	}
	t[key] = target
}

func (t _StrictTable[T]) Get(method, path string) (target T, found bool) {
	key := _StrictKey{
		method: method,
		path:   path,
	}
	target, found = t[key]
	return
}
