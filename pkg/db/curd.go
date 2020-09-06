package db

import (
	"context"
	"errors"
	"github.com/amit/file-download-manager/pkg/logger"
	"github.com/jinzhu/gorm"
	"time"
)

var DbContextKey = "DB_CONTEXT_KEY"

type Repoer interface {
	Find(ctx context.Context, models interface{}, id string, tableName string) error
	FindWhere(ctx context.Context, model interface{}, whereAttrs map[string]interface{}, tableName string) error
	Create(ctx context.Context, tableName string, model interface{}) error
	Update(ctx context.Context, model interface{}, args ...interface{}) error
	DBFromContext(context.Context) *gorm.DB
	Begin(context.Context) (context.Context, error)
	Commit(context.Context) error
	Rollback(context.Context) error
}

type Repo struct {
	DB *gorm.DB
}

var RepoClient Repoer

func (r *Repo) DBFromContext(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return r.DB
	}
	// If running transaction exists in context use that.
	contextDb, ok := ctx.Value(DbContextKey).(*gorm.DB)
	if ok {
		return contextDb
	}
	return r.DB
}

func (r *Repo) Find(ctx context.Context, models interface{}, id string, tableName string) error {
	if id == "" {
		logger.LogStatement("DB Find ERROR :: id is empty", nil)
		return errors.New("id is empty")
	}
	whereAttributes := map[string]interface{} {
		"id" : id,
	}
	err := r.DBFromContext(ctx).Table(tableName).Where(whereAttributes).First(models).Error
	if err != nil {
		logger.LogStatement("DB Find ERROR :: ", err)
		return err
	}
	return nil
}

func (r *Repo) FindWhere(ctx context.Context, model interface{}, whereAttrs map[string]interface{}, tableName string) error {
	err := r.DBFromContext(ctx).Table(tableName).Find(model, whereAttrs).Error
	if err != nil {
		logger.LogStatement("DB Find Where ERROR :: ", err)
	}
	return err
}

func (r *Repo) Create(ctx context.Context, tableName string, model interface{}) error {
	err := r.DBFromContext(ctx).Table(tableName).Create(model).Error
	if err != nil {
		logger.LogStatement("DB Create ERROR :: ", err)
	}
	return err
}

func (r *Repo) Update(ctx context.Context, model interface{}, attrs ...interface{}) error {
	updateAttrs := map[string]interface{}{
		"updated_at":	time.Now().Unix(),
	}
	attrs = append(attrs, updateAttrs)
	err := r.DBFromContext(ctx).Model(model).Update(attrs).Error
	if err != nil {
		logger.LogStatement("DB Update ERROR :: ", err)
		return err
	}
	return nil
}

func (r *Repo) Begin(ctx context.Context) (context.Context, error) {
	_, ok := ctx.Value(DbContextKey).(*gorm.DB)
	if ok {
		return ctx, nil
	}
	db := r.DB.Begin()
	err := db.Error
	if err != nil {
		logger.LogStatement("DB Transaction Begin ERROR :: ", err)
		return ctx, errors.New("internal server error")
	}
	return context.WithValue(ctx, DbContextKey, db), nil
}

func (r *Repo) Commit(ctx context.Context) error {
	err := r.DBFromContext(ctx).Commit().Error
	if err != nil {
		logger.LogStatement("DB Transaction Commit ERROR :: ", err)
		return errors.New("internal server error")
	}
	return err
}

func (r *Repo) Rollback(ctx context.Context) error {
	err := r.DBFromContext(ctx).Rollback().Error
	if err != nil {
		logger.LogStatement("DB Transaction Rollback ERROR :: ", err)
		return errors.New("internal server error")
	}
	return err
}