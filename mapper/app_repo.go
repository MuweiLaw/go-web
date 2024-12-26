package mapper

import (
	"database/sql"
	"go-web/entity/dao"
)

type AppRepository interface {
	Insert(app *dao.App) (int, error)
	Update(app dao.App) error
	DeleteByAppId(appId string) (bool, error)
	ExistsByAppId(appId string) (bool, error)
	FindByAppId(appId string) (*dao.App, error)
	ListAll() ([]dao.App, error)
	UpdateSecret(appId string, secret string) error
	UpdatePermissions(appId string, permissions string) error
}

func NewAppRepository(db *sql.DB) AppRepository {
	return &appRepo{db: db}
}

type appRepo struct {
	db *sql.DB
}

func (a appRepo) Insert(app *dao.App) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (a appRepo) Update(app dao.App) error {
	//TODO implement me
	panic("implement me")
}

func (a appRepo) DeleteByAppId(appId string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (a appRepo) ExistsByAppId(appId string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (a appRepo) FindByAppId(appId string) (*dao.App, error) {
	//TODO implement me
	panic("implement me")
}

func (a appRepo) ListAll() ([]dao.App, error) {
	//TODO implement me
	panic("implement me")
}

func (a appRepo) UpdateSecret(appId string, secret string) error {
	//TODO implement me
	panic("implement me")
}

func (a appRepo) UpdatePermissions(appId string, permissions string) error {
	//TODO implement me
	panic("implement me")
}
