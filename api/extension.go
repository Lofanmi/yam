package api

type Extension interface {
	Name() string
	Version() string
}
