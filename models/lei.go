package models

type Lei struct {
	Id   int    `json:"id" db:"id"`
	Nome string `json:"nome" db:"nome"`
}
