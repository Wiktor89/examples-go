package model

type Group struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Gr   Group  `json:"gr"`
}
