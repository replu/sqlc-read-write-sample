// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package sqlc

import (
	"time"
)

type User struct {
	ID        uint64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
