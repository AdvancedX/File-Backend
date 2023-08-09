package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"kratos-realworld/internal/biz"
	"strconv"
)

type User struct {
	ID           uint
	Email        string
	Username     string
	Bio          string
	Image        string
	PasswordHash string
}
type userRepo struct {
	data *Data
	log  *log.Helper
}
type profileRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
func (r *userRepo) CreateUser(ctx context.Context, u *biz.User) error {
	user := User{
		Email:        u.Email,
		Username:     u.Username,
		Bio:          u.Bio,
		Image:        u.Image,
		PasswordHash: u.PasswordHash,
	}
	collection := r.data.db.Collection("users")

	// 插入用户数据
	_, err := collection.InsertOne(ctx, user)
	return err
}
func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*biz.User, error) {
	// 创建过滤条件
	filter := bson.M{"email": email}

	// 获取用户集合对象
	collection := r.data.db.Collection("users")

	// 查询用户数据
	user := new(User)
	err := collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// 用户不存在
			return nil, errors.New(1, "user not found by email", "")
		}
		return nil, err
	}

	// 将查询到的用户数据转换为 biz.User 类型，并返回
	return &biz.User{
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		Bio:          user.Bio,
		Image:        user.Image,
		PasswordHash: user.PasswordHash,
	}, nil
}
func (r *userRepo) GetUserByID(ctx context.Context, id uint) (*biz.User, error) {
	// 将 uint 类型的 ID 转换为 string
	stringID := strconv.FormatUint(uint64(id), 10)

	// 将 string 类型的 ID 转换为 MongoDB 的 ObjectID
	objectID, err := primitive.ObjectIDFromHex(stringID)
	if err != nil {
		return nil, err
	}

	// 创建过滤条件
	filter := bson.M{"_id": objectID}

	// 获取用户集合对象
	collection := r.data.db.Collection("users")

	// 查询用户数据
	user := new(User)
	err = collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 用户不存在
			return nil, errors.New(2, "user not found by ID", "")
		}
		return nil, err
	}

	// 将查询到的用户数据转换为 biz.User 类型，并返回
	return &biz.User{
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		Bio:          user.Bio,
		Image:        user.Image,
		PasswordHash: user.PasswordHash,
	}, nil
}
func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (rv *biz.User, err error) {
	u := new(User)

	// Assuming you have a MongoDB session called "session"
	collection := r.data.db.Collection("users")

	// Create a filter for the "username" field
	filter := bson.M{"username": username}

	// Perform the find operation with the filter
	err = collection.FindOne(ctx, filter).Decode(u)
	if err != nil {
		return nil, err
	}

	return &biz.User{
		ID:           u.ID,
		Email:        u.Email,
		Username:     u.Username,
		Bio:          u.Bio,
		Image:        u.Image,
		PasswordHash: u.PasswordHash,
	}, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, in *biz.User) (rv *biz.User, err error) {
	// Create a filter for the "username" field
	filter := bson.M{"username": in.Username}

	// Create an update document with the fields to be updated
	update := bson.M{
		"$set": bson.M{
			"email":         in.Email,
			"bio":           in.Bio,
			"password_hash": in.PasswordHash,
			"image":         in.Image,
		},
	}

	// Assuming you have a MongoDB session called "session"
	collection := r.data.db.Collection("users")

	// Perform the update operation
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Fetch the updated document to return to the caller
	updatedUser := new(User)
	err = collection.FindOne(ctx, filter).Decode(updatedUser)
	if err != nil {
		return nil, err
	}

	return &biz.User{
		ID:           updatedUser.ID,
		Email:        updatedUser.Email,
		Username:     updatedUser.Username,
		Bio:          updatedUser.Bio,
		Image:        updatedUser.Image,
		PasswordHash: updatedUser.PasswordHash,
	}, nil
}
func (r *profileRepo) GetProfile(ctx context.Context, username string) (*biz.Profile, error) {
	collection := r.data.db.Collection("profiles")

	filter := bson.M{"username": username}
	var u biz.Profile
	err := collection.FindOne(ctx, filter).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 用户不存在的情况
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}
