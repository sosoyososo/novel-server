package DO

type BaseDO struct {
	AppType        int `json:"app"`
	SystemInfoDO   `comment:"系统信息"`
	BaseUserInfoDO `comment:"用户信息"`
}

type BaseUserInfoDO struct {
	UserID string `json:"userId" comment:"用户id"`
	OpenID string `json:"openId" comment:"用户openid"`
}

type SystemInfoDO struct {
	Device       string `json:"device" comment:"设备类型"`
	AppVersion   string `json:"appVersion" comment:"app版本"`
	WxVersion    string `json:"wxVersion" comment:"微信客户端版本"`
	WxSDKVersion string `json:"wxSdkVersion" comment:"微信SDK版本"`
	OsVersion    string `json:"osVersion" comment:"系统版本"`
	Os           string `json:"os" comment:"系统"`
}

type SearchDO struct {
	PageInfoDO `comment:"页面信息"`
	Key        string `json:"key" comment:"搜索关键字"`
}

type PageInfoDO struct {
	PageIndex  int `json:"page" comment:"分页页码 从0开始"`
	SizeOfPage int `json:"size" comment:"每页数据条数"`
}

func (p PageInfoDO) Page() int {
	if p.PageIndex >= 0 {
		return p.PageIndex
	}
	return 0
}

func (p PageInfoDO) Size() int {
	if p.SizeOfPage > 0 {
		if p.SizeOfPage > 50 {
			return 50
		}
		return p.SizeOfPage
	}
	return 20
}

type DetailDO struct {
	ID string `json:"id" comment:"id"`
}

type DetailPageDO struct {
	DetailDO
	PageInfoDO
}

type UserDO struct {
	UserID string `json:"userId" comment:"userId"`
}

type EmptyDO struct {
}

type AppBaseDO struct {
	AppName string `json:"appName"`
}
