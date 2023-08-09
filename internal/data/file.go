package data

import (
	"context"
	"kratos-realworld/internal/biz"
	"kratos-realworld/internal/pkg/utils"
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
	RelativePath *string            `bson:"videoRelativePath"`
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
		ID:          primitive.NewObjectID(),
		Type:        file.Type,
		Title:       file.Title,
		Description: file.Description,
		Tags:        file.Tags,
		UpdateTime:  &now,
	}
	_, err := v.collection.InsertOne(ctx, doc)
	return err
}

// Exists 判断文件是否存在
func (v *fileRepo) Exists(ctx context.Context, fileID string) (*biz.File, bool, error) {
	file := &File{}
	objectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return nil, false, err
	}
	err = v.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(file)
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
	bizVideos := make([]*biz.File, 0, len(files))
	for _, file := range files {
		bizVideos = append(bizVideos, &biz.File{
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
	return bizVideos, err
}

// DeleteOne 删除一个视频
func (v *fileRepo) DeleteOne(ctx context.Context, fileID string) error {
	idFromHex, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return err
	}
	_, err = v.collection.DeleteOne(ctx, bson.M{"_id": idFromHex})
	return err
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

func NewFileRepo(data *Data, logger log.Logger) biz.FileRepo {
	return &fileRepo{
		data:       data,
		collection: data.db.Collection(fileCollection),
		log:        log.NewHelper(log.With(logger, "module", "data/fileRepo")),
	}
}
