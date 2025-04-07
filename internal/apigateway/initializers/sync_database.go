package initializers

import (
	"github.com/KietAPCS/test_recruitment_assistant/internal/database/models"
)

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}