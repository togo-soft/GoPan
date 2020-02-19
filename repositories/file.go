package repositories

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server/models"
)

// FileRepo 文件操作仓库实体 实现了 FileRepoInterface 接口
type FileRepo struct {
}

// NewFileRepo 返回一个 FileRepo对象
func NewFileRepo() FileRepoInterface {
	return &FileRepo{}
}

func (this *FileRepo) UploadFile(username string, file *models.File) error {
	collection := mgo.Database("file").Collection(username)
	if _, err := collection.InsertOne(ctx, file); err != nil {
		return err
	}
	return nil
}

func (this *FileRepo) CreateDir(username,dirname string,file *models.File) error {
	collection := mgo.Database("file").Collection(username)
	if _, err := collection.InsertOne(ctx, file); err != nil {
		return err
	}
	return nil
}

func (this *FileRepo) DownloadFile(username string, id primitive.ObjectID) error {
	panic("implement me")
}

func (this *FileRepo) DeleteFile(username string, id primitive.ObjectID) error {
	collection := mgo.Database("file").Collection(username)
	_, err := collection.DeleteOne(ctx, bson.D{{"id", id}})
	return err
}

func (this *FileRepo) RenameFile(username, filename string, id primitive.ObjectID) error {
	collection := mgo.Database("file").Collection(username)
	filter := bson.D{{"id", id}}
	update := bson.D{{"$set", bson.D{{"filename", filename}}}}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (this *FileRepo) ShareFile(username string, id primitive.ObjectID) error {
	collection := mgo.Database("file").Collection(username)
	filter := bson.D{{"id", id}}
	update := bson.D{{"$set", bson.D{{"isshare", true}}}}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (this *FileRepo) DeleteDir(username string, id primitive.ObjectID) error {
	panic("implement me")
}

func (this *FileRepo) ListDir(username string, pid primitive.ObjectID) ([]models.File, error) {
	collection := mgo.Database("file").Collection(username)
	//排序规则
	var opts = options.Find().SetSort(bson.D{{"isdir",-1}})
	res, _ := collection.Find(ctx, bson.D{{"pid", pid}},opts)
	var result []models.File
	err := res.All(ctx, &result)
	return result, err
}

func (this *FileRepo) ListRoot(username string) ([]models.File, primitive.ObjectID, error) {
	collection := mgo.Database("file").Collection(username)
	var root = new(models.File)
	if err := collection.FindOne(ctx, bson.D{{"filename", "/"}}).Decode(root); err != nil {
		return nil, root.Id, err
	}
	// 获取子结构信息
	//排序规则
	var opts = options.Find().SetSort(bson.D{{"isdir",-1}})
	res, _ := collection.Find(ctx, bson.D{{"pid", root.Id}},opts)
	var result []models.File
	err := res.All(ctx, &result)
	return result, root.Id, err
}

func (this *FileRepo) FileInfo(username string, id primitive.ObjectID) (*models.File, error) {
	collection := mgo.Database("file").Collection(username)
	var result = new(models.File)
	err := collection.FindOne(ctx, bson.D{{"id", id}}).Decode(result)
	return result, err
}
