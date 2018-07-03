package server

type person struct {
	Id       int64    `json:"id"`
	Username string `json:"username"`
}

type message struct {
	From person `json:"from"`
	To   person `json:"to"`
	Body string `json:"body"`
}
