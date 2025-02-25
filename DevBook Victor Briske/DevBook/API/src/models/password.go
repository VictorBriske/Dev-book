package models

type Password struct {
	New    string `json:"new"`
	Actual string `json:"actual"`
}
