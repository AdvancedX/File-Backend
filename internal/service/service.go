package service

import (
	"context"
	"mime/multipart"

	"github.com/google/wire"

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
type CreateFileResponse struct{}
type UpdateFileResponse struct{}

func (b *BackendService) CreateFileHandler(ctx context.Context, req *CreateFileRequest) (*CreateFileResponse, error) {
	file := &biz.File{
		ID:          "",
		Type:        req.Type,
		Title:       req.Title,
		Description: req.Description,
		Tags:        req.Tags,
		FilePart:    req.FilePart,
	}
	err := b.fc.CreateFile(ctx, file)
	if err != nil {
		return nil, err
	}
	return &CreateFileResponse{}, nil
}
func (b *BackendService) UpdateFileHandler(ctx context.Context, req *UpdateFileRequest) (*UpdateFileResponse, error) {
	video := &biz.File{
		ID:          req.ID,
		Type:        req.Type,
		Title:       req.Title,
		Description: req.Description,
		Tags:        req.Tags,
		FilePart:    req.FilePart,
	}
	err := b.fc.UpdateFile(ctx, video)
	if err != nil {
		return nil, err
	}
	return &UpdateFileResponse{}, nil
}
