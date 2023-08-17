package data

import (
	"context"
	"kratos-realworld/internal/biz"
	"kratos-realworld/internal/pkg/utils"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const fileCollection = "file"

type File struct {
	ID           primitive.ObjectID `bson:"_id"`
	Type         string             `bson:"type,omitempty"`
	Title        string             `bson:"title,omitempty"`
	Description  string             `bson:"description,omitempty"`
	Tags         []string           `bson:"tags,omitempty"`
	UpdateTime   *time.Time         `bson:"updateTime"`
	RelativePath *string            `bson:"RelativePath"`
}

type fileRepo struct {
	data       *Data
	collection *mongo.Collection
	log        *log.Helper
}

// Save 保存文件记录
func (v *fileRepo) Save(ctx context.Context, file *biz.File) error {
	now := time.Now()
	doc := &File{
		ID:           primitive.NewObjectID(),
		Type:         file.Type,
		Title:        file.Title,
		Description:  file.Description,
		Tags:         file.Tags,
		UpdateTime:   &now,
		RelativePath: file.RelativePath,
	}
	_, err := v.collection.InsertOne(ctx, doc)
	return err
}

// Exists 判断文件是否存在
func (v *fileRepo) Exists(ctx context.Context, fileID string) (*biz.File, bool, error) {
	var file File
	objectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return nil, false, err
	}
	err = v.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&file)
	if err != nil {
		return nil, false, err
	}
	return &biz.File{
		ID:           file.ID.Hex(),
		Type:         file.Type,
		Title:        file.Title,
		Description:  file.Description,
		Tags:         file.Tags,
		FilePart:     nil,
		UpdateTime:   file.UpdateTime,
		RelativePath: file.RelativePath,
	}, !file.ID.IsZero(), nil
}

// Update 更新文件记录
func (v *fileRepo) Update(ctx context.Context, file *biz.File) error {
	now := time.Now()
	hex, err := primitive.ObjectIDFromHex(file.ID)
	if err != nil {
		return err
	}
	doc := &File{
		ID:           hex,
		Type:         file.Type,
		Title:        file.Title,
		Description:  file.Description,
		Tags:         file.Tags,
		UpdateTime:   &now,
		RelativePath: file.RelativePath,
	}
	_, err = v.collection.UpdateByID(ctx, hex, bson.D{{Key: "$set", Value: doc}})
	return err
}

// ListByType 按类型返回文件列表
func (v *fileRepo) ListByType(ctx context.Context, fileType string) ([]*biz.File, error) {
	opts := options.Find().SetSort(bson.D{{Key: "updateTime", Value: -1}})
	cursor, err := v.collection.Find(ctx, bson.M{"type": fileType}, opts)
	if err != nil {
		return nil, err
	}
	var files []*File
	err = cursor.All(ctx, &files)
	if err != nil {
		return nil, err
	}
	bizFiles := make([]*biz.File, 0, len(files))
	for _, file := range files {
		bizFiles = append(bizFiles, &biz.File{
			ID:           file.ID.Hex(),
			Type:         file.Type,
			Title:        file.Title,
			Description:  file.Description,
			Tags:         file.Tags,
			FilePart:     nil,
			UpdateTime:   file.UpdateTime,
			RelativePath: file.RelativePath,
		})
	}
	return bizFiles, err
}

// DeleteOne 删除一个文件
func (v *fileRepo) DeleteOne(ctx context.Context, fileID string) error {
	idFromHex, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return err
	}
	var document File
	err = v.collection.FindOne(ctx, bson.M{"_id": idFromHex}).Decode(&document)
	filepath := "file/" + *document.RelativePath
	filename := "/" + document.Title + document.Type
	dirpath := strings.ReplaceAll(filepath, filename, "")
	_, err = v.collection.DeleteOne(ctx, bson.M{"_id": idFromHex})
	if err != nil {
		return err
	}
	err = os.RemoveAll(dirpath)
	if err != nil {
		return err
	}
	return nil
}

// ListTagsByType 按类型返回视频标签列表
func (v *fileRepo) ListTagsByType(ctx context.Context, fileType string) ([]string, error) {
	opts := options.Find().SetSort(bson.D{{Key: "updateTime", Value: -1}})
	cursor, err := v.collection.Find(ctx, bson.M{"type": fileType}, opts)
	if err != nil {
		return nil, err
	}
	var files []*File
	err = cursor.All(ctx, &files)
	if err != nil {
		return nil, err
	}
	var tags []string
	for _, file := range files {
		for _, tag := range file.Tags {
			if !utils.SliceContainsAny(tags, tag) {
				tags = append(tags, tag)
			}
		}
	}
	return tags, err
}
func (v *fileRepo) FindPathByID(ctx context.Context, fileID string) (*string, error) {
	cursor, err := v.collection.Find(ctx, bson.M{"fileID": fileID})
	if err != nil {
		return nil, err
	}
	var files = File{RelativePath: nil}
	cursor.
		All(ctx, &files)
	return files.RelativePath, nil
}
func (v *fileRepo) DownloadFile(ctx context.Context, fileID string) (*biz.DownloadFileReply, error) {
	idFromHex, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return nil, err
	}
	var document File
	err = v.collection.FindOne(ctx, bson.M{"_id": idFromHex}).Decode(&document)
	if err != nil {
		return nil, err
	}
	filePath := path.Join("file", *document.RelativePath)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return &biz.DownloadFileReply{
		Title:    document.Title + document.Type,
		FilePart: file,
	}, nil
}
func NewFileRepo(data *Data, logger log.Logger) biz.FileRepo {
	return &fileRepo{
		data:       data,
		collection: data.db.Collection(fileCollection),
		log:        log.NewHelper(log.With(logger, "module", "data/fileRepo")),
	}
}
