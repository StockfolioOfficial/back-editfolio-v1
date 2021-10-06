package repository

import (
	"context"
	"database/sql"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"github.com/stockfolioofficial/back-editfolio/util/gormx"
	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	db.AutoMigrate(&domain.User{})
	return &repo{
		db: db,
	}
}


type repo struct {
	db *gorm.DB
}

func (r *repo) Get() *gorm.DB {
	return r.db
}

func (r *repo) Save(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *repo) Transaction(ctx context.Context, fn func(userRepo domain.UserTxRepository) error, options ...*sql.TxOptions) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&repo{db: tx})
	}, options...)
}

func (r *repo) With(tx gormx.Tx) domain.UserTxRepository {
	return &repo{db: tx.Get()}
}
