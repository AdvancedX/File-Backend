package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport/http"
	"io"
	"kratos-realworld/internal/pkg/middleware/auth"
	"kratos-realworld/internal/pkg/utils"
	"kratos-realworld/internal/service"
	"path"
	"strings"
)

var files = []string{".txt", ".doc", ".pdf", ".xlsx", ".pptx", ".jpg", ".png", ".zip", ".tar", ".gz"}

func CreateFileHandler(backend *service.BackendService) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		req := &service.CreateFileRequest{}
		err := ctx.Request().ParseMultipartForm(20 << 26)
		if err != nil {
			return err
		}
		req.Type = ctx.Request().MultipartForm.Value["type"][0]
		req.Title = ctx.Request().MultipartForm.Value["title"][0]
		req.Description = ctx.Request().MultipartForm.Value["description"][0]
		req.Tags = ctx.Request().MultipartForm.Value["tags"]
		req.FilePart = ctx.Request().MultipartForm.File["FilePart"][0]
		if !utils.SliceContainsAny(files, strings.ToLower(path.Ext(req.FilePart.Filename))) {
			return errors.New("文件格式错误，请输入其中的一种")
		}
		http.SetOperation(ctx, auth.OperationBackendCustomCreateVideo)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return backend.CreateFileHandler(ctx, req.(*service.CreateFileRequest))
		})
		out, err := h(ctx, req)
		if err != nil {
			return err
		}
		reply, _ := out.(*service.CreateFileResponse)
		return ctx.Result(200, reply)
	}
}

func UpdateFileHandler(backend *service.BackendService) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		req := &service.UpdateFileRequest{}
		err := ctx.Request().ParseMultipartForm(20 << 26)
		if err != nil {
			return err
		}
		req.ID = ctx.Request().MultipartForm.Value["id"][0]
		if err != nil {
			return err
		}
		req.Type = ctx.Request().MultipartForm.Value["type"][0]
		req.Title = ctx.Request().MultipartForm.Value["title"][0]
		req.Description = ctx.Request().MultipartForm.Value["description"][0]
		req.Tags = strings.Split(ctx.Request().MultipartForm.Value["tags"][0], ",")
		FilePartParam := ctx.Request().MultipartForm.File["FilePart"]
		if len(FilePartParam) != 0 {
			if !utils.SliceContainsAny(files, strings.ToLower(path.Ext(FilePartParam[0].Filename))) {
				fmt.Println("文件格式错误，请输入其中的一种")
				return err
			}
			req.FilePart = FilePartParam[0]
		}
		http.SetOperation(ctx, auth.OperationBackendCustomUpdateVideo)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return backend.UpdateFileHandler(ctx, req.(*service.UpdateFileRequest))
		})
		out, err := h(ctx, req)
		if err != nil {
			return err
		}

		reply, _ := out.(*service.UpdateFileResponse)
		return ctx.Result(200, reply)
	}
}
func DownloadFileHandler(backend *service.BackendService) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		req := &service.DownloadFileRequest{}
		err := ctx.Request().ParseMultipartForm(20 << 26)
		if err != nil {
			return err
		}
		req.ID = ctx.Request().MultipartForm.Value["id"][0]
		if err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return backend.DownloadFileHandler(ctx, req.(*service.DownloadFileRequest))
		})

		out, err := h(ctx, req)
		if err != nil {
			return err
		}

		reply, _ := out.(*service.DownloadFileReply)
		err = streamFile(reply, ctx.Response())
		return nil
	}
}

// streamFile 流式文件传输
func streamFile(reply *service.DownloadFileReply, w http.ResponseWriter) error {
	w.Header().Set("Content-Disposition", "attachment;filename="+reply.Title)

	if _, err := io.Copy(w, reply.FilePart); err != nil {
		// 处理错误
		return err
	}
	defer reply.FilePart.Close()
	return nil
}
