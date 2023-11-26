package di

// kind represents a type without using reflection.
type kind[T any] struct{}

func intoKind[T any]() kind[T] {
    return kind[T]{}
}
