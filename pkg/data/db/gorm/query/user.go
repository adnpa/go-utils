package query

import (
	"adnpa/go-utils/pkg/data/db/gorm/db"
	"adnpa/go-utils/pkg/data/db/gorm/model"
	"context"
	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &UserDao{db.NewDBClient(ctx)}
}

func (dao *UserDao) FindUserByUserName(userName string) (r *model.User, err error) {
	err = dao.Model(&model.User{}).
		Where("name = ?", userName).Find(&r).Error
	if err != nil {
		return
	}

	return
}

// CreateUser 创建用户
func (dao *UserDao) CreateUser(in *model.User) (err error) {
	return dao.Model(&model.User{}).Create(&in).Error
}
