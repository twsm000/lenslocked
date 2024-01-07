package result

func MustGet[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func ExtractValue[A any, B any](a A, b B) A {
	return a
}
