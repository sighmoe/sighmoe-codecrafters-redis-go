package main

var s map[string]string

func StoreInit() {
	s = make(map[string]string)
}

func Set(key string, val string) {
	s[key] = val
}

func Get(key string) (string, bool) {
	val, ok := s[key]
	return val, ok
}
