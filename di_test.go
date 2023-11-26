package di

import (
    "github.com/jProgr/assert"
    "testing"
)

type SomeStruct struct {
    field int
}

func NewSomeStruct(value int) (*SomeStruct, error) {
    return &SomeStruct{value}, nil
}

func TestBindsTypeToProvider(test *testing.T) {
    container := NewContainer()
    provider := func(_ *Container) (*SomeStruct, error) {
        return NewSomeStruct(100)
    }

    Bind(container, provider)

    untypedProvider := container.bindings[intoKind[*SomeStruct]()]
    boundProvider := untypedProvider.(Provider[*SomeStruct])
    result, _ := boundProvider(container)
    expected, _ := provider(container)

    assert.Equals(test, expected.field, result.field)
}

func TestBindsNamedToProvider(test *testing.T) {
    id := "some_id"
    container := NewContainer()
    provider := func(_ *Container) (*SomeStruct, error) {
        return NewSomeStruct(100)
    }

    BindNamed(container, id, provider)

    untypedProvider := container.namedBindings[id]
    boundProvider := untypedProvider.(Provider[*SomeStruct])
    result, _ := boundProvider(container)

    assert.Equals(test, 100, result.field)
}

func TestBindsIfNotAlreadyRegistered(test *testing.T) {
    container := NewContainer()
    provider := func(_ *Container) (*SomeStruct, error) {
        return NewSomeStruct(100)
    }

    BindIf(container, provider)
    result, err := Make[*SomeStruct](container)
    assert.Nil(test, err)
    assert.Equals(test, 100, result.field)

    BindIf(container, func(_ *Container) (*SomeStruct, error) {
        return NewSomeStruct(200)
    })
    result, err = Make[*SomeStruct](container)
    assert.Nil(test, err)
    assert.Equals(test, 100, result.field)
}

func TestResolvesADependency(test *testing.T) {
    container := NewContainer()
    provider := func(_ *Container) (*SomeStruct, error) {
        return NewSomeStruct(100)
    }
    Bind(container, provider)

    instance, err := Make[*SomeStruct](container)

    assert.Nil(test, err)
    assert.Equals(test, 100, instance.field)
}

func TestResolvesANamedDependency(test *testing.T) {
    id := "some_id"
    container := NewContainer()
    provider := func(_ *Container) (*SomeStruct, error) {
        return NewSomeStruct(100)
    }
    BindNamed(container, id, provider)

    instance, err := MakeNamed[*SomeStruct](container, id)

    assert.Nil(test, err)
    assert.Equals(test, 100, instance.field)
}
