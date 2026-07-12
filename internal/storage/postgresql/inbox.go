package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alexgul25/notify-svc/internal/domain"
	"github.com/alexgul25/notify-svc/internal/inbox"
)

type InboxStorage struct {
	db *sql.DB
}

func NewInboxStorage(db *sql.DB) *InboxStorage {
	return &InboxStorage{db: db}
}

func (is *InboxStorage) InsertRecord(ctx context.Context, record inbox.Record) error {
	const op = "postgresql.InboxStorage.InsertRecord"

	query := `
		INSERT INTO inbox (id, topic, process_status, payload, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO NOTHING
	`

	result, err := is.db.ExecContext(ctx, query, record.ID, record.Topic, record.ProcessStatus, record.Payload, record.CreatedAt)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if n == 0 {
		return domain.ErrMsgDoubleSend
	}

	return nil
}

func (is *InboxStorage) SelectPending(ctx context.Context, topic string, limit int) ([]inbox.Record, error) {
	const op = "postgresql.InboxStorage.SelectPending"

	query := `
		SELECT id, payload, created_at
		FROM inbox
		WHERE topic = $1 AND process_status = $2
		ORDER BY created_at
		LIMIT $3
	`

	rows, err := is.db.QueryContext(ctx, query, topic, inbox.StatusPending, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var records []inbox.Record
	for rows.Next() {
		record := inbox.Record{Topic: topic, ProcessStatus: inbox.StatusPending}
		err := rows.Scan(
			&record.ID,
			&record.Payload,
			&record.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return records, nil
}

func (is *InboxStorage) MarkAsProcessed(ctx context.Context, id string) error {
	const op = "postgresql.InboxStorage.MarkAsProcessed"

	query := `
		UPDATE inbox
		SET process_status = $1
		WHERE id = $2
	`

	result, err := is.db.ExecContext(ctx, query, inbox.StatusProcessed, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("record %s not found", id)
	}

	return nil
}
