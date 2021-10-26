package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"github.com/stockfolioofficial/back-editfolio/util/gormx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewCustomerRepository(db *gorm.DB) domain.CustomerRepository {
	db.AutoMigrate(&domain.Customer{})
	return &repo{db: db}
}

type repo struct {
	db *gorm.DB
}

func (r *repo) GetById(ctx context.Context, userId uuid.UUID) (customer *domain.Customer, err error) {
	var entity domain.Customer

	err = r.db.WithContext(ctx).First(&entity, userId).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}

	customer = &entity
	return
}

func (r *repo) GetByUsername(ctx context.Context, username string) (user *domain.User, err error) {
	var entity domain.User
	err = r.db.WithContext(ctx).
		Where("`username` = ?", username).
		First(&entity).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}

	user = &entity
	return
}

func (r *repo) Get() *gorm.DB {
	return r.db
}

func (r *repo) Save(ctx context.Context, customer *domain.Customer) error {
	//TODO
	//refactor
	//gormx.Upsert(ctx, r.db, customer)
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(customer).Error
}

func (r *repo) With(tx gormx.Tx) domain.CustomerTxRepository {
	return &repo{db: tx.Get()}
}
