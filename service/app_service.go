package service

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-web/entity/dao"
	"go-web/mapper"
	"strings"
)

type AppService interface {
	//InitAppSecret(appList []entity.App) error
	Create(req dao.CreateAppRequest) (*dao.App, error)
	//Remove(appId string) error
	//Find(appId string) (*entity.App, error)
	//ListAll() ([]dto.AppListItem, error)
	//ResetSecret(appId string) error
	//ListPermissions(appId string) ([]dto.AppPermissionDTO, error)
	//UpdatePermissions(appId string, permissions []string) error
	//MqttLogin(appId string, secret string) (*dto.AppMqList, error)
}

func NewAppService(db *sql.DB) AppService {
	return &appService{
		appRepo: mapper.NewAppRepository(db),
		//permRepo: mapper.NewPermissionRepository(db),
	}
}

type appService struct {
	appRepo mapper.AppRepository
	//permRepo      mapper.PermissionRepository
	//permGroupRepo mapper.PermissionGroupRepository
}

func (s *appService) Create(req dao.CreateAppRequest) (*dao.App, error) {
	appId := strings.ToLower(req.AppId)
	exists, err := s.appRepo.ExistsByAppId(appId)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(fmt.Sprintf("app %s already exists", appId))
	}

	if len(strings.TrimSpace(req.Secret)) == 0 {
		// 生成一个32字节长的Secret字符串
		secret, err := generateSecret(32)
		if err != nil {
			log.Warning("Error generating secret:", err)
		}
		log.Info("Generated Secret:", secret)
	}

	app := &dao.App{
		AppId:  appId,
		Secret: req.Secret,
	}
	id, err := s.appRepo.Insert(app)
	if err != nil {
		return nil, err
	}
	app.Id = id
	return app, nil
}

// generateSecret 生成指定长度的随机Secret字符串
func generateSecret(length int) (string, error) {
	// 创建一个字节切片来存储随机字节
	bytes := make([]byte, length)

	// 使用crypto/rand包生成随机字节
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// 将字节切片编码为Base64字符串
	secret := base64.URLEncoding.EncodeToString(bytes)

	return secret, nil
}
