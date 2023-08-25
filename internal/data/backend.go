package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"kratos-realworld/internal/biz"
)

type backendRepo struct {
	data *Data
	log  *log.Helper
}

// NewBackendRepo .
func NewBackendRepo(data *Data, logger log.Logger) biz.BackendRepo {
	return &backendRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
func NewTransaction(d *Data) biz.Transaction {
	return d
}

func (r *backendRepo) Save(ctx context.Context, g *biz.Backend) (*biz.Backend, error) {
	return g, nil
}

func (r *backendRepo) Update(ctx context.Context, g *biz.Backend) (*biz.Backend, error) {
	return g, nil
}

func (r *backendRepo) FindByID(context.Context, int64) (*biz.Backend, error) {
	return nil, nil
}

func (r *backendRepo) ListByHello(context.Context, string) ([]*biz.Backend, error) {
	return nil, nil
}

func (r *backendRepo) ListAll(context.Context) ([]*biz.Backend, error) {
	return nil, nil
}
