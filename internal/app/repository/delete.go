package repository

import (
	"database/sql"
	"fmt"
	"system_components/internal/app/ds"
)

func (r *Repository) DeleteComponentRawSQL(id int) error {
	query := "UPDATE system_components SET status = $1 WHERE id = $2"
	res := r.db.Exec(query, string(ds.StatusDeleted), id)
	if res.Error != nil {
		return fmt.Errorf("ошибка SQL UPDATE: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}