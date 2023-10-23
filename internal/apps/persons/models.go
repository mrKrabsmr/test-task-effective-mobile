package persons

import "github.com/google/uuid"

type Person struct {
	ID         uuid.UUID `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	Surname    string    `db:"surname" json:"surname"`
	Patronymic string    `db:"patronymic" json:"patronymic"`
	Age        int       `db:"age" json:"age"`
	Gender     string    `db:"gender" json:"gender"`
	Nation     string    `db:"nation" json:"nation"`
}
