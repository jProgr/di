package di

// Provider is a function that returns instances of dependencies.
type Provider[T any] func(*Container) (T, error)

// ProviderId identifies named providers.
type ProviderId = string

// Container holds the mappings of identifiers and their providers.
type Container struct {
    bindings           map[any]any
    namedBindings      map[ProviderId]any
    singletons         map[any]any
    resolvedSingletons map[any]any
    failedSingletons   map[any]error
}

func NewContainer() *Container {
    return &Container{
        bindings:           make(map[any]any),
        namedBindings:      make(map[ProviderId]any),
        singletons:         make(map[any]any),
        resolvedSingletons: make(map[any]any),
        failedSingletons:   make(map[any]error),
    }
}

// Bind marks the function provider as the dependency resolver for type T.
func Bind[T any](container *Container, provider Provider[T]) {
    container.bindings[intoKind[T]()] = provider
}

// Bind marks the function provider as the dependency resolver for type T
// only if there is already none registered for that type.
func BindIf[T any](container *Container, provider Provider[T]) {
    if isRegistered(intoKind[T](), container) {
        return
    }

    Bind(container, provider)
}

// Bind marks the function provider as the dependency resolver for id.
func BindNamed[T any](container *Container, id ProviderId, provider Provider[T]) {
    container.namedBindings[id] = provider
}

// BindSingleton marks the function provider as the dependency resolver for type T
// as singleton.
func BindSingleton[T any](container *Container, provider Provider[T]) {
    container.singletons[intoKind[T]()] = provider
}

// Make resolves a dependency of type T.
func Make[T any](container *Container) (T, error) {
    instance, err := resolveBinding[T](container)

    return instance, err
}

// MakeNamed resolves a dependency by ProviderId.
func MakeNamed[T any](container *Container, id ProviderId) (T, error) {
    if !isRegisteredAsNamed(id, container) {
        return intoZero[T](), ProviderNotFoundError
    }

    return resolveProvider[T](container.namedBindings[id], container)
}

// resolveBinding looks for a provider that satisfies T, calls it and
// returns the instance. Only normal and singleton bindings are considered.
func resolveBinding[T any](container *Container) (T, error) {
    kind := intoKind[T]()

    // No provider found.
    if !isRegistered(kind, container) {
        return intoZero[T](), ProviderNotFoundError
    }

    // Search for a singleton binding.
    if isRegisteredAsSingleton(kind, container) {
        // Was previously resolved with error.
        if err, isResolvedWithError := container.failedSingletons[kind]; isResolvedWithError {
            return intoZero[T](), err
        }

        // Was previously resolved with success.
        if untypedInstance, isResolved := container.resolvedSingletons[kind]; isResolved {
            instance := untypedInstance.(T)

            return instance, nil
        }

        // First time resolving.
        if instance, err := resolveProvider[T](container.singletons[kind], container); err != nil {
            container.failedSingletons[kind] = err
        } else {
            container.resolvedSingletons[kind] = instance
        }

        return resolveBinding[T](container)
    }

    // Search for a normal binding.
    if isRegisteredAsNormal(kind, container) {
        return resolveProvider[T](container.bindings[kind], container)
    }

    // This shouldn't happen. It's either registered or not.
    return intoZero[T](), ContainerError
}

func resolveProvider[T any](untypedProvider any, container *Container) (T, error) {
    actualProvider := untypedProvider.(Provider[T])

    return actualProvider(container)
}

func isRegistered[T any](kind kind[T], container *Container) bool {
    return isRegisteredAsNormal(kind, container) || isRegisteredAsSingleton(kind, container)
}

func isRegisteredAsNormal[T any](kind kind[T], container *Container) bool {
    _, isRegistered := container.bindings[kind]

    return isRegistered
}

func isRegisteredAsSingleton[T any](kind kind[T], container *Container) bool {
    _, isRegistered := container.singletons[kind]

    return isRegistered
}

func isRegisteredAsNamed(id ProviderId, container *Container) bool {
    _, isRegistered := container.namedBindings[id]

    return isRegistered
}

func intoZero[T any]() (zero T) {
    return
}
