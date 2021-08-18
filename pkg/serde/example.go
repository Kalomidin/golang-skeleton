//go:generate easyjson -all -omit_empty $GOFILE
package serde

type ExampleModel struct {
	ID  string            `json:"ID"`
	Foo string            `json:"foo"`
	Bar map[string]string `json:"bar"`
}

type ExampleOutput struct {
	ID     string            `json:"ID"`
	FooBar map[string]string `json:"foo_bar"`
}
