package repositories

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"server/models"
	"server/utils"
)

// FileRepo 文件操作仓库实体 实现了 FileRepoInterface 接口
type FileRepo struct {
}

// NewFileRepo 返回一个 FileRepo对象
func NewFileRepo() FileRepoInterface {
	return &FileRepo{}
}

// UploadFile 上传文件
func (this *FileRepo) UploadFile(username string, file *models.File) error {
	collection := mgo.Database("file").Collection(username)
	//插入记录
	if _, err := collection.InsertOne(ctx, file); err != nil {
		return err
	}
	//修改存储占用率
	var sto = new(models.FileStorage)
	err := collection.FindOne(ctx, bson.D{{"username", username}}).Decode(sto)
	if err != nil {
		return err
	}
	filter := bson.D{{"username", username}}
	update := bson.D{{"$set", bson.D{{"usedsize", file.Size + sto.UsedSize}}}}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

// CreateDir 创建目录
func (this *FileRepo) CreateDir(username, dirname string, file *models.File) error {
	collection := mgo.Database("file").Collection(username)
	if _, err := collection.InsertOne(ctx, file); err != nil {
		return err
	}
	return nil
}

// DeleteFile 删除文件
func (this *FileRepo) DeleteFile(username string, id primitive.ObjectID) error {
	collection := mgo.Database("file").Collection(username)
	var file = new(models.File)
	if err := collection.FindOne(ctx, bson.D{{"id", id}}).Decode(file); err != nil {
		return err
	}
	if _, err := collection.DeleteOne(ctx, bson.D{{"id", id}}); err != nil {
		return err
	}
	//修改存储占用率
	var sto = new(models.FileStorage)
	err := collection.FindOne(ctx, bson.D{{"username", username}}).Decode(sto)
	if err != nil {
		return err
	}
	filter := bson.D{{"username", username}}
	update := bson.D{{"$set", bson.D{{"usedsize", sto.UsedSize - file.Size}}}}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

// MoveFile 移动文件位置
func (this *FileRepo) MoveFile(username string, id, pid primitive.ObjectID) error {
	collection := mgo.Database("file").Collection(username)
	filter := bson.D{{"id", id}}
	update := bson.D{{"$set", bson.D{{"pid", pid}}}}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

// RenameFile 文件/目录重命名
func (this *FileRepo) RenameFile(username, filename string, id primitive.ObjectID) error {
	collection := mgo.Database("file").Collection(username)
	filter := bson.D{{"id", id}}
	update := bson.D{{"$set", bson.D{{"filename", filename}}}}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

// OTTHShareFile 第三方获取文件共享列表
func (this *FileRepo) OTTHShareFile(username string, id primitive.ObjectID) ([]models.File, error) {
	collection := mgo.Database("file").Collection(username)
	var opts = options.Find().SetSort(bson.D{{"isdir", -1}})
	res, _ := collection.Find(ctx, bson.D{{"pid", id}}, opts)
	var result []models.File
	err := res.All(ctx, &result)
	return result, err
}

// ShareList 共享列表
func (this *FileRepo) ShareList(username string) ([]models.File, error) {
	collection := mgo.Database("file").Collection(username)
	//排序规则
	var opts = options.Find().SetSort(bson.D{{"isdir", -1}})
	res, _ := collection.Find(ctx, bson.D{{"isshare", true}}, opts)
	var result []models.File
	err := res.All(ctx, &result)
	return result, err
}

// ShareFile 共享文件
func (this *FileRepo) ShareFile(username string, id primitive.ObjectID, fsk string) error {
	collection := mgo.Database("file").Collection(username)
	filter := bson.D{{"id", id}}
	update := bson.D{{"$set", bson.D{{"isshare", true}, {"fsk", fsk}}}}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

// CancelShare 取消共享
func (this *FileRepo) CancelShare(username string, id primitive.ObjectID) error {
	collection := mgo.Database("file").Collection(username)
	filter := bson.D{{"id", id}}
	update := bson.D{{"$set", bson.D{{"isshare", false}, {"fsk", ""}}}}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

// getFileIDList 删除文件操作前 先获取所有子项的ID
func getFileIDList(id primitive.ObjectID, list map[primitive.ObjectID]string, collection *mongo.Collection) map[primitive.ObjectID]string {
	var opts = options.Find().SetSort(bson.D{{"isdir", 1}})
	res, _ := collection.Find(ctx, bson.D{{"pid", id}}, opts)
	var result []models.File
	_ = res.All(ctx, &result)
	if len(result) == 0 {
		//空目录 返回即可
		return nil
	}
	for _, v := range result {
		list[v.Id] = v.HashCode
		if v.IsDir {
			//文件夹 执行子集查询
			getFileIDList(v.Id, list, collection)
		}
	}
	return list
}

// deleteFileFromDFS 从DFS删除文件
func deleteFileFromDFS(list map[primitive.ObjectID]string) {
	server := utils.GetConfig().File.DFS
	for _, value := range list {
		if value != "" {
			//只删除文件
			log.Println(server + "/delete?md5=" + value)
			_, _ = http.Get(server + "/delete?md5=" + value)
		}
	}
}

// deleteFileFromMongoDB 从mongodb删除文件
func deleteFileFromMongoDB(username string, list map[primitive.ObjectID]string, collection *mongo.Collection) error {
	var usedSize float64
	var file = new(models.File)
	for key, _ := range list {
		//查询信息
		if err := collection.FindOne(ctx, bson.D{{"id", key}}).Decode(file); err != nil {
			return err
		}
		//统计此次删除恢复的空间大小
		usedSize += file.Size
		//删除文件
		if _, err := collection.DeleteOne(ctx, bson.D{{"id", key}}); err != nil {
			return err
		}
	}
	//修改存储占用率
	var sto = new(models.FileStorage)
	if err := collection.FindOne(ctx, bson.D{{"username", username}}).Decode(sto); err != nil {
		return err
	}
	filter := bson.D{{"username", username}}
	update := bson.D{{"$set", bson.D{{"usedsize", sto.UsedSize - usedSize}}}}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

// DeleteDir 删除文件夹操作
func (this *FileRepo) DeleteDir(username string, id primitive.ObjectID) error {
	collection := mgo.Database("file").Collection(username)
	//将要删除的目标节点ID放入map中
	list := make(map[primitive.ObjectID]string)
	list[id] = ""
	//递归将下一级ID也放入map
	log.Println(getFileIDList(id, list, collection))
	//执行删除操作
	deleteFileFromDFS(list)
	return deleteFileFromMongoDB(username, list, collection)
}

// ListDir 罗列指定目录的信息
func (this *FileRepo) ListDir(username string, pid primitive.ObjectID) ([]models.File, error) {
	collection := mgo.Database("file").Collection(username)
	//排序规则
	var opts = options.Find().SetSort(bson.D{{"isdir", -1}})
	res, _ := collection.Find(ctx, bson.D{{"pid", pid}}, opts)
	var result []models.File
	err := res.All(ctx, &result)
	return result, err
}

// ListRoot 罗列文件列表
func (this *FileRepo) ListRoot(username string) ([]models.File, primitive.ObjectID, error) {
	collection := mgo.Database("file").Collection(username)
	var root = new(models.File)
	if err := collection.FindOne(ctx, bson.D{{"filename", "/"}}).Decode(root); err != nil {
		return nil, root.Id, err
	}
	// 获取子结构信息
	//排序规则
	var opts = options.Find().SetSort(bson.D{{"isdir", -1}})
	res, _ := collection.Find(ctx, bson.D{{"pid", root.Id}}, opts)
	var result []models.File
	err := res.All(ctx, &result)
	return result, root.Id, err
}

// ListSecret 私密文件列表
func (this *FileRepo) ListSecret(username string) ([]models.File, primitive.ObjectID, error) {
	collection := mgo.Database("file").Collection(username)
	var root = new(models.File)
	if err := collection.FindOne(ctx, bson.D{{"filename", "/#"}}).Decode(root); err != nil {
		return nil, root.Id, err
	}
	// 获取子结构信息
	//排序规则
	var opts = options.Find().SetSort(bson.D{{"isdir", -1}})
	res, _ := collection.Find(ctx, bson.D{{"pid", root.Id}}, opts)
	var result []models.File
	err := res.All(ctx, &result)
	return result, root.Id, err
}

// FileInfo 文件信息
func (this *FileRepo) FileInfo(username string, id primitive.ObjectID) (*models.File, error) {
	collection := mgo.Database("file").Collection(username)
	var result = new(models.File)
	err := collection.FindOne(ctx, bson.D{{"id", id}}).Decode(result)
	return result, err
}

// ShareFileInfo 共享文件信息
func (this *FileRepo) ShareFileInfo(username string, fsk string) (*models.File, error) {
	collection := mgo.Database("file").Collection(username)
	var result = new(models.File)
	err := collection.FindOne(ctx, bson.D{{"fsk", fsk}}).Decode(result)
	return result, err
}

// UpdateFileStorage 修改用户存储统计总大小
func (this *FileRepo) UpdateFileStorage(username string, totalSize float64) error {
	collection := mgo.Database("file").Collection(username)
	filter := bson.D{{"username", username}}
	update := bson.D{{"$set", bson.D{{"totalsize", totalSize}}}}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

// CheckRatio 检测使用量是否大于总量
func (this *FileRepo) CheckRatio(username string) error {
	collection := mgo.Database("file").Collection(username)
	var result = new(models.FileStorage)
	err := collection.FindOne(ctx, bson.D{{"username", username}}).Decode(result)
	if err != nil {
		return err
	}
	//比较大小 UsedSize TotalSize
	if result.UsedSize > result.TotalSize {
		return errors.New("空间已满,无法继续上传")
	}
	return nil
}

// UsageRate 返回用户使用磁盘比率
func (this *FileRepo) UsageRate(username string) (*models.FileStorage, error) {
	collection := mgo.Database("file").Collection(username)
	var result = new(models.FileStorage)
	err := collection.FindOne(ctx, bson.D{{"username", username}}).Decode(result)
	return result, err
}

// CollectionList 收藏文件列表
func (this *FileRepo) CollectionList(username string) ([]models.FileCollection, primitive.ObjectID, error) {
	collection := mgo.Database("file").Collection(username)
	var root = new(models.File)
	if err := collection.FindOne(ctx, bson.D{{"filename", "/@"}}).Decode(root); err != nil {
		return nil, root.Id, err
	}
	// 获取子结构信息
	var opts = options.Find().SetSort(bson.D{{"collectiontime", -1}})
	res, _ := collection.Find(ctx, bson.D{{"pid", root.Id}}, opts)
	var result []models.FileCollection
	err := res.All(ctx, &result)
	return result, root.Id, err
}

// CollectionFile 收藏文件
func (this *FileRepo) CollectionFile(username string, fc *models.FileCollection) error {
	collection := mgo.Database("file").Collection(username)
	if _, err := collection.InsertOne(ctx, fc); err != nil {
		return err
	}
	return nil
}

// CancelCollection 取消收藏文件
func (this *FileRepo) CancelCollection(username string, id primitive.ObjectID) error {
	collection := mgo.Database("file").Collection(username)
	if _, err := collection.DeleteOne(ctx, bson.D{{"id", id}}); err != nil {
		return err
	}
	return nil
}

// FindFileByFilename 通过文件名 查询文件实体
func (this *FileRepo) FindFileByFilename(username, filename string) (*models.File, error) {
	collection := mgo.Database("file").Collection(username)
	var result = new(models.File)
	err := collection.FindOne(ctx, bson.D{{"filename", filename}}).Decode(result)
	return result, err
}
