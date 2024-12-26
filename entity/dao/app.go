package dao

type App struct {
	Id     int    `json:"id" db:"id"`
	AppId  string `json:"appId" db:"app_id"`
	Secret string `json:"secret" db:"secret"`
	//Permissions     string   `json:"permissions" db:"permissions"`
	MqPubs string `json:"mq_pubs" db:"mq_pubs"`
	MqSubs string `json:"mq_subs" db:"mq_subs"`
	//PermissionArray []string `json:"permissionArray"`
	//MqPubArray      []string `json:"MqPubArray"`
	//MqSubArray      []string `json:"MqSubArray"`
	IsEnable    string `json:"isEnable" db:"is_enable"`
	LockVersion string `json:"lockVersion" db:"lock_version"`
	Creator     string `json:"creator" db:"creator"`
	Updater     string `json:"updater" db:"updater"`
	CreateTime  string `json:"createTime" db:"create_time"`
	UpdateTime  string `json:"updateTime" db:"update_time"`
}

type AppListItem struct {
	Id     int    `json:"id"`
	AppId  string `json:"appId"`
	Secret string `json:"secret"`
}

type CreateAppRequest struct {
	AppId       string   `json:"appId"`
	Secret      string   `json:"secret"`
	Permissions []string `json:"permissions"`
}

type UpdateAppPermissionsRequest struct {
	Permissions []string `json:"permissions"`
}

type AppMqList struct {
	Pubs []string `json:"pubs"`
	Subs []string `json:"subs"`
}
