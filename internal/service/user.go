package service

import (
	"context"

	v1 "kratos-realworld/api/backend/v1"
	"kratos-realworld/internal/biz"
)

func (b *BackendService) Login(ctx context.Context, req *v1.LoginRequest) (reply *v1.UserReply, err error) {
	rv, err := b.uc.Login(ctx, req.User.Email, req.User.Password)
	if err != nil {
		return nil, err
	}
	return &v1.UserReply{
		User: &v1.UserReply_User{
			Email:    rv.Email,
			Username: rv.Username,
			Token:    rv.Token,
		},
	}, nil
}

func (b *BackendService) Register(ctx context.Context, req *v1.RegisterRequest) (reply *v1.UserReply, err error) {
	u, err := b.uc.Register(ctx, req.User.Username, req.User.Email, req.User.Password)
	if err != nil {
		return nil, err
	}
	return &v1.UserReply{
		User: &v1.UserReply_User{
			Email:    u.Email,
			Username: u.Username,
			Token:    u.Token,
		},
	}, nil
}
func (b *BackendService) GetCurrentUser(ctx context.Context, req *v1.GetCurrentUserRequest) (reply *v1.UserReply, err error) {
	u, err := b.uc.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.UserReply{
		User: &v1.UserReply_User{
			Username: u.Username,
		},
	}, nil
}
func (b *BackendService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (rep *v1.UserReply, err error) {
	u, err := b.uc.UpdateUser(ctx, &biz.UserUpdate{
		Email:    req.User.GetEmail(),
		Username: req.User.GetUsername(),
		Password: req.User.GetPassword(),
	})
	if err != nil {
		return nil, err
	}
	return &v1.UserReply{
		User: &v1.UserReply_User{
			Username: u.Username,
			Email:    u.Email,
		},
	}, nil
}
