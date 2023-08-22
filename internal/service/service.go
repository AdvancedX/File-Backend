package service

import (
	"context"
	"github.com/google/wire"
	v1 "kratos-realworld/api/backend/v1"
	"mime/multipart"
	"os"

	"kratos-realworld/internal/biz"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewBackendService)

type CreateFileRequest struct {
	Type        string
	Title       string
	Content     string
	Description string
	Tags        []string
	FilePart    *multipart.FileHeader
}
type UpdateFileRequest struct {
	ID          string
	Type        string
	Title       string
	Description string
	Tags        []string
	FilePart    *multipart.FileHeader
}
type CreateFileResponse struct {
	reply string
}
type UpdateFileResponse struct{}
type DownloadFileRequest struct {
	ID string
}

type DownloadFileReply struct {
	Title    string
	FilePart *os.File
}

func (b *BackendService) CreateFileHandler(ctx context.Context, req *CreateFileRequest) (*CreateFileResponse, error) {
	file := &biz.File{
		ID:          "",
		Type:        req.Type,
		Title:       req.Title,
		Description: req.Description,
		Tags:        req.Tags,
		FilePart:    req.FilePart,
	} //Relative Part nil,taglist only one
	err := b.fc.CreateFile(ctx, file)
	if err != nil {
		return nil, err
	}
	return &CreateFileResponse{}, nil
}
func (b *BackendService) UpdateFileHandler(ctx context.Context, req *UpdateFileRequest) (*UpdateFileResponse, error) {
	file := &biz.File{
		ID:          req.ID,
		Type:        req.Type,
		Title:       req.Title,
		Description: req.Description,
		Tags:        req.Tags,
		FilePart:    req.FilePart,
	}
	err := b.fc.UpdateFile(ctx, file)
	if err != nil {
		return nil, err
	}
	return &UpdateFileResponse{}, nil
}
func (b *BackendService) ListFileByType(ctx context.Context, in *v1.ListFileRequest) (*v1.ListFileReply, error) {
	files, err := b.fc.ListByType(ctx, in.FileType)
	if err != nil {
		return nil, err
	}
	results := make([]*v1.File, 0, len(files))
	for _, file := range files {
		result := &v1.File{
			Id:          file.ID,
			Type:        file.Type,
			Title:       file.Title,
			Description: file.Description,
			Tags:        file.Tags,
			UpdateTime:  file.UpdateTime.Local().Format("2006-01-02 15:04:05"),
		}
		if file.RelativePath != nil {
			result.FilePath = *file.RelativePath
		}
		results = append(results, result)
	}
	return &v1.ListFileReply{Files: results}, nil
}
func (b *BackendService) DeleteFile(ctx context.Context, in *v1.DeleteFileRequest) (*v1.DeleteFileReply, error) {
	return &v1.DeleteFileReply{}, b.fc.DeleteOne(ctx, in.FileID)
}
func (b *BackendService) FindFileByName(ctx context.Context, in *v1.FindFileRequest) (*v1.FindFileReply, error) {
	file, err := b.fc.FindByName(ctx, in.FileName)
	if err != nil {
		return nil, err
	}
	files := &v1.File{
		Id:          file.ID,
		Type:        file.Type,
		Title:       file.Title,
		Description: file.Description,
		Tags:        file.Tags,
		FilePath:    *file.RelativePath,
		UpdateTime:  file.UpdateTime.Local().Format("2006-01-02 15:04:05"),
	}
	return &v1.FindFileReply{File: files}, nil
}

func (b *BackendService) DownloadFileHandler(ctx context.Context, req *DownloadFileRequest) (*DownloadFileReply, error) {
	result, err := b.fc.DownloadFile(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &DownloadFileReply{
		Title:    result.Title,
		FilePart: result.FilePart,
	}, nil
}
