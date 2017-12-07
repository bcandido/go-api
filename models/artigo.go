package models

type Artigo struct {
	Numero    string `json:"numero" db:"numero"`
	Descrisao string `json:"descrisao" db:"descrisao"`
	Lei       string `json:"lei" db:"lei"`
}
