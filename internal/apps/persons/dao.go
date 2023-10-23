package persons

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	core "github.com/mrKrabsmr/commerce-edu-api/internal/apps"
)

type dao struct {
	*sqlx.DB
}

func newDAO() *dao {
	return &dao{
		core.GetDB(),
	}
}

const (
	selectAllFields = "SELECT * FROM persons"
)

func (d *dao) getSearchAndFilter(filterObject *Person, searchValue string, p *core.Paginate) ([]*Person, error) {
	var persons []*Person
	var args []any

	query := selectAllFields

	if searchValue != "" {
		searchValue = fmt.Sprintf("%%%s%%", searchValue)
		query += fmt.Sprintf(` WHERE name ILIKE ($%d) OR
					surname ILIKE ($%d) OR
					patronymic ILIKE ($%d)`, len(args)+1, len(args)+2, len(args)+3)

		args = append(args, searchValue, searchValue, searchValue)
	}

	if filterObject != nil {
		if searchValue != "" {
			query += " AND "
		} else {
			query += " WHERE"
		}

		if filterObject.Name != "" {
			query += fmt.Sprintf(" name = $%d", len(args)+1)
			args = append(args, filterObject.Name)
		}

		if filterObject.Surname != "" {
			query += fmt.Sprintf(" surname = $%d", len(args)+1)
			args = append(args, filterObject.Surname)
		}

		if filterObject.Patronymic != "" {
			query += fmt.Sprintf(" patronymic = $%d", len(args)+1)
			args = append(args, filterObject.Patronymic)
		}

		if filterObject.Age != -1 {
			query += fmt.Sprintf(" age = $%d", len(args)+1)
			args = append(args, filterObject.Age)
		}

		if filterObject.Gender != "" {
			query += fmt.Sprintf(" gender = $%d", len(args)+1)
			args = append(args, filterObject.Gender)
		}

		if filterObject.Nation != "" {
			query += fmt.Sprintf(" nation = $%d", len(args)+1)
			args = append(args, filterObject.Nation)
		}
	}

	query += " ORDER BY name, surname, patronymic"

	if p != nil {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)

		args = append(args, p.Limit, p.StartWith)
	}

	if err := d.Select(&persons, query, args...); err != nil {
		return nil, err
	}

	return persons, nil
}

func (d *dao) getOne(id uuid.UUID) (*Person, error) {
	var person Person

	query := selectAllFields + " WHERE id = $1"

	if err := d.Get(&person, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &person, nil
}

func (d *dao) delete(id uuid.UUID) error {
	query := "DELETE FROM persons WHERE id = $1"

	if _, err := d.Exec(query, id); err != nil {
		return err
	}

	return nil
}

func (d *dao) create(person *Person) error {
	query := `INSERT INTO persons
				VALUES($1, $2, $3, $4, $5, $6, $7)`

	if _, err := d.Exec(
		query, person.ID, person.Name, person.Surname, person.Patronymic,
		person.Age, person.Gender, person.Nation,
	); err != nil {
		return err
	}

	return nil
}

func (d *dao) update(person *Person) error {
	query := `UPDATE persons SET 
				name = $1,
				surname = $2,
				patronymic = $3,
				age = $4,
				gender = $5,
				nation = $6
				WHERE id = $7`

	if _, err := d.Exec(
		query, person.Name, person.Surname, person.Patronymic,
		person.Age, person.Gender, person.Nation, person.ID,
	); err != nil {
		return err
	}

	return nil
}
