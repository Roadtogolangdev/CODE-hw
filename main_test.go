package storage_test

import (
	"code-hw/internal/storage"
	"context"

	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddNote(t *testing.T) {
	// Создаем mock-объекты для базы данных
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Ожидаем запрос для добавления заметки
	query := `INSERT INTO notes \(text, user_id\) VALUES \(\$1, \$2\) RETURNING id, text, user_id, created_at`
	mock.ExpectQuery(query).
		WithArgs("Test note", 123).
		WillReturnRows(sqlmock.NewRows([]string{"id", "text", "user_id", "created_at"}).
			AddRow(1, "Test note", 123, time.Now()))

	// Создаем экземпляр Storage с mock-базой данных
	store := storage.NewSqlStorage(db)

	// Выполняем тестируемую функцию
	ctx := context.Background()
	note, err := store.AddNote(ctx, "Test note", 123)

	// Проверяем, что ошибок не было
	assert.NoError(t, err)

	// Проверяем, что заметка была добавлена корректно
	assert.Equal(t, 1, note.ID)
	assert.Equal(t, "Test note", note.Text)

	// Проверяем, что все ожидания sqlmock выполнены
	assert.NoError(t, mock.ExpectationsWereMet())
}
