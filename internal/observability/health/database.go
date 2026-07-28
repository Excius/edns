package health

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseChecker struct {
	db *pgxpool.Pool
}

func NewDatabaseChecker(db *pgxpool.Pool) *DatabaseChecker {
	return &DatabaseChecker{
		db: db,
	}
}

func (d *DatabaseChecker) Check(ctx context.Context) CheckResult {
	if err := d.db.Ping(ctx); err != nil {
		return CheckResult{
			Name:   "database",
			Status: "down",
			Error:  err.Error(),
		}
	}

	return CheckResult{
		Name:   "database",
		Status: "up",
	}
}
