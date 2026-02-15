package config

type Provider[T any] func() *T
