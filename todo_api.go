package main

type TodoJSON struct {
	Id          int32  `json:"id"`
	Created     int32  `json:"created"`
	Status      string `json:"status"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}
