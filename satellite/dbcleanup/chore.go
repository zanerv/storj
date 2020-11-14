// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package dbcleanup

import (
	"context"
	"time"

	"github.com/spacemonkeygo/monkit/v3"
	"github.com/zeebo/errs"
	"go.uber.org/zap"

	"storj.io/common/sync2"
	"storj.io/storj/satellite/orders"
)

var (
	// Error the default dbcleanup errs class.
	Error = errs.Class("dbcleanup error")

	mon = monkit.Package()
)

// Config defines configuration struct for dbcleanup chore.
type Config struct {
	SerialsInterval           time.Duration `help:"how often to delete expired serial numbers" default:"4h"`
	SerialsDeleteEnabled      bool          `help:"xxx" default:"true"`
	SerialsDeleteTimeout      int64         `help:"xxx" default:"5000"`
	SerialsDeleteTimeoutCount int           `help:"xxx" default:"5"`
	SerialsDeleteLimit        int64         `help:"xxx" default:"1000"`
}

// Chore for deleting DB entries that are no longer needed.
//
// architecture: Chore
type Chore struct {
	log    *zap.Logger
	orders orders.DB

	Serials *sync2.Cycle
	config  Config
}

// NewChore creates new chore for deleting DB entries.
func NewChore(log *zap.Logger, orders orders.DB, config Config) *Chore {
	return &Chore{
		log:    log,
		orders: orders,

		config: config,
	}
}

// Run starts the db cleanup chore.
func (chore *Chore) Run(ctx context.Context) (err error) {
	defer mon.Task()(&ctx)(&err)
	return chore.Serials.Run(ctx, chore.deleteExpiredSerials)
}

func (chore *Chore) deleteExpiredSerials(ctx context.Context) (err error) {
	defer mon.Task()(&ctx)(&err)
	chore.log.Debug("deleting expired serial numbers")

	now := time.Now()

	var options *orders.SerialDeleteOptions
	if chore.config.SerialsDeleteEnabled {
		options = &orders.SerialDeleteOptions{
			Timeout:      chore.config.SerialsDeleteTimeout,
			TimeoutCount: chore.config.SerialsDeleteTimeoutCount,
			Limit:        chore.config.SerialsDeleteLimit,
		}
	}

	deleted, err := chore.orders.DeleteExpiredSerials(ctx, now, options)
	if err != nil {
		chore.log.Error("deleting expired serial numbers", zap.Error(err))
	} else {
		chore.log.Debug("expired serials deleted", zap.Int("items deleted", deleted))
	}

	deleted, err = chore.orders.DeleteExpiredConsumedSerials(ctx, now)
	if err != nil {
		chore.log.Error("deleting expired serial numbers", zap.Error(err))
	} else {
		chore.log.Debug("expired serials deleted", zap.Int("items deleted", deleted))
	}

	return nil
}

// Close stops the dbcleanup chore.
func (chore *Chore) Close() error {
	chore.Serials.Close()
	return nil
}
