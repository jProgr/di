package di

import (
    "errors"
)

var ProviderNotFoundError error = errors.New("Provider not found")
var ContainerError error = errors.New("It is not possible for a type to be on a state different from registration or not registration")
