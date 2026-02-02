package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Krokozabra213/schools_backend/services/sso/domain"
)

func (r *PostgresRepository) UpdateUser(ctx context.Context, params domain.UpdateUser) error {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if params.Username != nil {
		setParts = append(setParts, fmt.Sprintf("username = $%d", argIndex))
		args = append(args, *params.Username)
		argIndex++
	}

	if params.Email != nil {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, *params.Email)
		argIndex++
	}

	if params.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *params.Name)
		argIndex++
	}

	if params.Surname != nil {
		setParts = append(setParts, fmt.Sprintf("surname = $%d", argIndex))
		args = append(args, *params.Surname)
		argIndex++
	}

	if params.IsMale != nil {
		setParts = append(setParts, fmt.Sprintf("is_male = $%d", argIndex))
		args = append(args, *params.IsMale)
		argIndex++
	}

	// Нечего обновлять
	if len(setParts) == 0 {
		return nil
	}

	// Добавляем ID
	args = append(args, params.ID)

	query := fmt.Sprintf(`
        UPDATE users
        SET %s
        WHERE id = $%d AND deleted_at IS NULL
    `, strings.Join(setParts, ", "), argIndex)

	result, err := r.DB.Exec(ctx, query, args...)
	if err != nil {
		return r.handleError(err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
