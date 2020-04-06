package repositories

import "server/models"

// UserRepo 操作用户模型 实现了UserRepoInterface 接口
type LogRepo struct {
}

// NewUserRepo 返回一个UserRepo对象
func NewLogRepo() LogRepoInterface {
	return &LogRepo{}
}

func (this *LogRepo) AddLog(log *models.UserLog) (int64, error) {
	return engine.Insert(log)
}

func (this *LogRepo) UpdateLog(log *models.UserLog) (int64, error) {
	return engine.Id(log.UserId).Update(log)
}

func (this *LogRepo) QueryLog(id int64) (*models.UserLog, error) {
	var log = &models.UserLog{}
	if _, err := engine.Where("uid = ?", id).Get(log); err != nil {
		return nil, err
	}
	return log, nil
}

func (this *LogRepo) LogList() ([]models.UserAndLog, error) {
	all := make([]models.UserAndLog, 0)
	err := engine.Table(&models.UserLog{}).Join("LEFT", "user", "user_log.user_id = user.id").Find(&all)
	return all, err
}
