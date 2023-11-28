# di

di implements a container for dependency injection for Go without using reflection.

## How does it not use reflection?

The problem is mapping types to functions. Some common solutions are to call `reflect.TypeOf()` or to use `fmt` with `%T` (which also uses reflection). Instead, this package creates empty structs with a generic. This captures the type and can be used for the mappings.

```golang
type kind[T any] struct{} // Usable as key for a map.
```

The downside is that there is not that much that can be done with that. It is not possible to get a string representation, know with which type is being dealt with inside the package, implement methods that deal with the providers or dependencies in the container due to being done through generics.

## Usage

Create a new container:

```golang
container := di.NewContainer()
```

Register providers:

```golang
provider := func(_ *di.Container) (*SomeStruct, error) {
    return NewSomeStruct(100)
}

di.Bind(container, provider)
```

Resolve dependencies:

```golang
instance, err = di.Make[*SomeStruct](container)
```

### Binding singletons

In this case, the provider will be registered and only called the first time the instance is needed. The following calls will return the same instance.

Register:

```golang
provider := func(_ *di.Container) (*SomeStruct, error) {
    return NewSomeStruct(100)
}

di.BindSingleton(container, provider)
```

Resolve:

```golang
instance, err = di.Make[*SomeStruct](container)
```

### Binding named providers

Optionally, the container also supports named dependencies by passing a string. This way, it is possible to have more than provider for the same type:

```golang
id := "some_id"
container := di.NewContainer()
provider := func(_ *di.Container) (*SomeStruct, error) {
    return NewSomeStruct(100)
}
di.BindNamed(container, id, provider)

instance, err := di.MakeNamed[*SomeStruct](container, id)
```
