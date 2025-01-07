package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-web/entity/dao"
	"go-web/mapper"
	"strings"
	"sync"
)

type AppService interface {
	Create(req dao.CreateAppRequest) (*dao.App, error)
}

var (
	as   AppService
	once sync.Once
)

type appService struct {
	appRepo mapper.AppRepository
}

func GetAppService() *AppService {
	if as == nil {
		once.Do(func() {
			as = &appService{*mapper.GetAppRepo()}
		})
	}
	return &as
}

// Create 创建一个APP
func (s *appService) Create(req dao.CreateAppRequest) (*dao.App, error) {
	var appid = strings.ToLower(req.AppId)
	var secret string

	exists, err := s.ExistsByAppId(appid)
	if err != nil {
		log.Warning("Error generating secret:", err)
		return nil, err
	}
	if exists {
		return nil, errors.New(fmt.Sprintf("app <%s> already exists", appid))
	}

	if len(strings.TrimSpace(req.Secret)) == 0 {
		// 生成一个32字节长的Secret字符串
		secret, err = generateSecret(32)
		if err != nil {
			log.Warning("Error generating secret:", err)
		}
		log.Infof("Appid<%s> generated Secret: %s", appid, secret)
	}

	app := &dao.App{
		AppId:  appid,
		Secret: secret,
	}
	id, err := s.appRepo.Insert(app)
	if err != nil {
		return nil, err
	}
	app.Id = id
	return app, nil
}

// ExistsByAppId 判断应用是否存在
func (s *appService) ExistsByAppId(appid string) (bool, error) {
	exists, err := s.appRepo.ExistsByAppId(appid)
	if err != nil {
		log.Warnf("Exists by appid: %s, Error: %s", appid, err.Error())
		return false, err
	}
	log.Infof("Appid <%s> exists: %v", appid, exists)
	return exists, nil
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
