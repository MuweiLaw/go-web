package mapper

import (
	"database/sql"
	"errors"
	log "github.com/sirupsen/logrus"
	"go-web/db/mysql"
	"go-web/entity/dao"
	"sync"
	"time"
)

type AppRepository interface {
	Insert(app *dao.App) (int64, error)
	Update(app dao.App) error
	DeleteByAppId(appId string) (bool, error)
	ExistsByAppId(appId string) (bool, error)
	FindByAppId(appId string) (*dao.App, error)
	ListAll() ([]dao.App, error)
	UpdateSecret(appId string, secret string) error
	UpdatePermissions(appId string, permissions string) error
}

var (
	ar   AppRepository
	once sync.Once
)

func GetAppRepo() *AppRepository {
	if ar == nil {
		once.Do(func() {
			ar = &appRepo{db: mysql.GetMysql()}
		})
	}
	return &ar
}

type appRepo struct {
	db *sql.DB
}

func (ar appRepo) getMysql() *sql.DB {
	return ar.db
}

func (ar appRepo) Insert(app *dao.App) (int64, error) {
	db := ar.getMysql()
	// 获取当前时间
	currentTime := time.Now()

	// 格式化时间为字符串
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	res, err := db.Exec("INSERT INTO `auth_app`  (app_id, secret, is_enable, creator, updater, create_time, update_time) VALUES (?,?,?,?,?,?,?)",
		app.AppId, app.Secret, 1, "admin", "admin", currentTime, formattedTime,
	)
	if err != nil {
		log.Printf("Appid：<%s> 新增失败!", app.AppId)
		return -1, err
	}
	//log.Infof("新增APP执行结果%v", &res)
	id, err := res.LastInsertId()
	if err != nil {
		return id, err
	}
	if id < 1 {
		return id, errors.New("新增失败！")
	}
	return id, err
}

func (ar appRepo) Update(app dao.App) error {
	//TODO implement me
	panic("implement me")
}

func (ar appRepo) DeleteByAppId(appId string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

// ExistsByAppId 查询数据库判断appid是否存在
func (ar appRepo) ExistsByAppId(appid string) (bool, error) {
	db := ar.getMysql()
	rows, err := db.Query("SELECT EXISTS(SELECT 1 FROM `auth_app` WHERE `app_id` = ?) AS appid_exists", appid)
	if err != nil {
		log.Printf("Appid：<%s> 查询失败!", appid)
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Printf(" rows close err: %s", err.Error())
		}
	}(rows)

	// 查询数据
	var appidExists bool

	for rows.Next() {
		if err := rows.Scan(&appidExists); err != nil {
			log.Errorf("Appid: %s 查询结果映射异常, err: %v", appid, err.Error())
			return false, err
		}
		log.Infof("Appid: %s, Exists: %t", appid, appidExists)
	}

	// 检查迭代是否因为错误而提前结束
	if err := rows.Err(); err != nil {
		log.Errorf("因错误：<%v>而提前结束, Appid: %s, Rows: %v", err.Error(), appid, rows)
		return false, err
	}
	return appidExists, err
}

func (ar appRepo) FindByAppId(appId string) (*dao.App, error) {
	//TODO implement me
	panic("implement me")
}

func (ar appRepo) ListAll() ([]dao.App, error) {
	//TODO implement me
	panic("implement me")
}

func (ar appRepo) UpdateSecret(appId string, secret string) error {
	//TODO implement me
	panic("implement me")
}

func (ar appRepo) UpdatePermissions(appId string, permissions string) error {
	//TODO implement me
	panic("implement me")
}
