package model

type BaseRes[T any] struct {
	Msg  string `json:"msg"`
	Data *T     `json:"data"`
}
