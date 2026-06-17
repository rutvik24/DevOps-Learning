package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/rutvik/todo-backend/internal/models"
	"gorm.io/gorm"
)

func All() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "20240617000001_create_todos_table",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Todo{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&models.Todo{})
			},
		},
	}
}

func Up(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, All())
	return m.Migrate()
}

func Down(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, All())
	return m.RollbackLast()
}

func DownAll(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, All())
	return m.RollbackTo("0")
}
