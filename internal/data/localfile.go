/******************************************************************************
* FILENAME:      file.go
*
* AUTHORS:       Xie Rongwang START DATE: 周六 11月 19 2022
*
* LAST MODIFIED: 星期六, 十一月 19th 2022, 上午10:06
*
* CONTACT:       rongwang.xie@smartmore.com
******************************************************************************/

package data

import (
	"io"
	"mime/multipart"
	"os"
	"path"

	"github.com/go-kratos/kratos/v2/log"
	"kratos-realworld/internal/biz"
	"kratos-realworld/internal/conf"
)

type filelocalRepo struct {
	data *Data
	conf *conf.Data_File

	log *log.Helper
}

func NewFileLocalRepo(data *Data, conf *conf.Data, logger log.Logger) biz.FileLocalRepo {
	return &filelocalRepo{
		data: data,
		conf: conf.File,
		log:  log.NewHelper(log.With(logger, "module", "data/file")),
	}
}

func (f *filelocalRepo) SaveLocalFile(relativePath string, localfile *multipart.FileHeader) error {
	fileOpen, err := localfile.Open()
	if err != nil {
		return err
	}
	defer fileOpen.Close()
	savePath := path.Join(f.conf.BasePath, relativePath)
	err = os.MkdirAll(path.Dir(savePath), os.ModePerm)
	if err != nil {
		return err
	}
	save, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer save.Close()
	_, err = io.Copy(save, fileOpen)
	if err != nil {
		return err
	}
	return nil
}
