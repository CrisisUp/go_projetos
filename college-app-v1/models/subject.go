// models/subject.go
package models

// Subject representa uma matéria na universidade.
type Subject struct {
	ID      string `json:"id"`      // ID único da matéria (ex: "BSI101")
	Name    string `json:"name"`    // Nome da matéria (ex: "Programação Orientada a Objetos")
	Year    int    `json:"year"`    // Ano em que a matéria é oferecida (ex: 1, 2, 3, 4)
	Credits int    `json:"credits"` // Créditos da matéria (ex: 4)
}
