package repositories

import (
	"server/models"
)

// UserRepo 操作用户模型 实现了UserRepoInterface 接口
type GroupRepo struct {
}

// NewUserRepo 返回一个UserRepo对象
func NewGroupRepo() GroupRepoInterface {
	return &GroupRepo{}
}

func (this *GroupRepo) AddGroup(group *models.Group) (int64, error) {
	return engine.Insert(group)
}

func (this *GroupRepo) UpdateGroup(group *models.Group) (int64, error) {
	return engine.Id(group.Id).Update(group)
}

func (this *GroupRepo) DeleteGroup(id int64) (int64, error) {
	return engine.Delete(&models.Group{Id: id})
}

func (this *GroupRepo) GroupList() ([]*models.Group, error) {
	all := make([]*models.Group, 0)
	err := engine.Find(&all)
	return all, err
}
