package biz

import (
	"context"
	"os"
	"path"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"

	"kratos-realworld/internal/conf"
)

type FileUsecase struct {
	localfile FileLocalRepo
	file      FileRepo
	conf      *conf.Data
	tm        Transaction

	log *log.Helper
}
type ProfileUsecase struct {
	pr ProfileRepo

	log *log.Helper
}
type DownloadFileReply struct {
	Title    string
	FilePart *os.File
}

func (v *FileUsecase) CreateFile(ctx context.Context, file *File) error {
	err := v.file.AvoidRepeatedFile(ctx, file.Title, file.Type)
	if err != nil {
		return err
	}
	intermediatePath := uuid.New().String()
	fileRelativePath := path.Join(v.conf.File.FilePath, intermediatePath, file.FilePart.Filename)
	// 单个文件，串行上传
	err = v.localfile.SaveLocalFile(fileRelativePath, file.FilePart)
	if err != nil {
		return err
	}
	file.RelativePath = &fileRelativePath
	return v.file.Save(ctx, file)
}
func (v *FileUsecase) UpdateFile(ctx context.Context, file *File) error {
	exist, ok, err := v.file.Exists(ctx, file.ID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New(500, "file not exists", "文件不存在")
	}
	file.RelativePath = exist.RelativePath

	intermediatePath := uuid.New().String()
	if file.FilePart != nil {
		fileRelativePath := path.Join(v.conf.File.FilePath, intermediatePath, file.FilePart.Filename)
		// 单个视频文件，串行上传
		err = v.localfile.SaveLocalFile(fileRelativePath, file.FilePart)
		if err != nil {
			return err
		}
		file.RelativePath = &fileRelativePath
	}
	return v.file.Update(ctx, file)
}

func (v *FileUsecase) ListByType(ctx context.Context, fileType string) ([]*File, error) {
	return v.file.ListByType(ctx, fileType)
}
func (v *FileUsecase) FindByName(ctx context.Context, FileName string) (*File, error) {
	return v.file.FindFileByName(ctx, FileName)
}
func (v *FileUsecase) DeleteOne(ctx context.Context, fileID string) error {
	return v.file.DeleteOne(ctx, fileID)
}
func (v *FileUsecase) DownloadFile(ctx context.Context, fileID string) (*DownloadFileReply, error) {
	return v.file.DownloadFile(ctx, fileID)
}

func NewFileUsecase(localfile FileLocalRepo, file FileRepo, conf *conf.Data, tm Transaction, logger log.Logger) *FileUsecase {
	return &FileUsecase{
		localfile: localfile,
		file:      file,
		conf:      conf,
		tm:        tm,
		log:       log.NewHelper(log.With(logger, "module", "biz/file")),
	}
}
