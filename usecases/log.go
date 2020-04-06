package usecases

import (
	"net/http"
	"server/models"
	"server/repositories"
	"server/utils"
	"time"
)

// ulr 是仓库存储层的一个实例
var ulr = repositories.NewLogRepo()

func AddLog(uid int64, r *http.Request, operate string) (int64, error) {
	var log = &models.UserLog{
		Id:               0,
		UserId:           uid,
		LastTime:         time.Time{},
		CurrentTime:      time.Now(),
		LastIP:           "",
		CurrentIP:        utils.RemoteIp(r),
		LastAddress:      "",
		CurrentAddress:   utils.IpGeolocation(utils.RemoteIp(r)),
		LastUserAgent:    "",
		CurrentUserAgent: r.Header.Get("User-Agent"),
		LastOperate:      "",
		CurrentOperate:   operate,
	}
	//写入数据库
	return ulr.AddLog(log)
}

func UpdateLog(uid int64, r *http.Request, operate string) (int64, error) {
	var last, _ = ulr.QueryLog(uid)
	var current = &models.UserLog{
		Id:               0,
		UserId:           uid,
		LastTime:         last.CurrentTime,
		CurrentTime:      time.Now(),
		LastIP:           last.CurrentIP,
		CurrentIP:        utils.RemoteIp(r),
		LastAddress:      last.CurrentAddress,
		CurrentAddress:   utils.IpGeolocation(utils.RemoteIp(r)),
		LastUserAgent:    last.CurrentUserAgent,
		CurrentUserAgent: r.Header.Get("User-Agent"),
		LastOperate:      last.CurrentOperate,
		CurrentOperate:   operate,
	}
	//写入数据库
	return ulr.UpdateLog(current)
}
