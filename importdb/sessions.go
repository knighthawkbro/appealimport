package importdb

import (
	"github.com/jinzhu/gorm"
)

// Session (Public) -
type Session struct {
	gorm.Model

	UUID   string
	Email  string
	UserID string
}
