package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"kratos-realworld/internal/conf"
	"path"
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

func (v *FileUsecase) CreateFile(ctx context.Context, file *File) error {
	intermediatePath := uuid.New().String()
	fileRelativePath := path.Join(v.conf.File.FilePath, intermediatePath, file.FilePart.Filename)
	// 单个文件，串行上传
	err := v.localfile.SaveLocalFile(fileRelativePath, file.FilePart)
	if err != nil {
		return err
	}
	file.RelativePath = &fileRelativePath
	return v.file.Save(ctx, file)
}
func (v *FileUsecase) UpdateFile(ctx context.Context, file *File) error {
	return v.tm.ExecTx(ctx, func(ctx context.Context) error {
		exist, ok, err := v.file.Exists(ctx, file.ID)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New(500, "video not exists", "文件不存在")
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
	})
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
