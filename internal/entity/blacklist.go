package entity

import "time"

type BlacklistedUser struct {
	ID        int       `db:"id"`
	Phone     string    `db:"phone"`
	Username  string    `db:"username"`
	Cause     string    `db:"cause"`
	Timestamp time.Time `db:"timestamp"`
	CausedBy  string    `db:"caused_by"`
}
