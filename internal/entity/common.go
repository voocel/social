package entity

type Request struct {
	Cmd    string `json:"cmd"`
	Params struct {
		Kind     string `json:"kind"`
		Sender   int    `json:"sender"`
		Receiver int    `json:"receiver"`
		Content  string `json:"content"`
	} `json:"params"`
}
