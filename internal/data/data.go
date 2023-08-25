package data

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"kratos-realworld/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewUserRepo, NewBackendRepo, NewFileRepo, NewFileLocalRepo, NewTransaction)

// Data .
type Data struct {
	client *mongo.Client
	err    error
	db     *mongo.Database
	test   *mongo.Collection
}

func (d *Data) ExecTx(ctx context.Context, f func(ctx context.Context) error) error {
	session, err := d.db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		return nil, f(ctx)
	})
	return err
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, db *mongo.Database) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db: db}, cleanup, nil
}

func NewDB(c *conf.Data) (*mongo.Database, error) {
	var (
		client *mongo.Client
		err    error
		db     *mongo.Database
	)
	// 连接MongoDB
	if client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(c.Database.Dsn).SetConnectTimeout(5*time.Second)); err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// 检查是否存在名为 "kratos" 的数据库
	databaseNames, err := client.ListDatabaseNames(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to list database names: %w", err)
	}

	kratosDatabaseExists := false
	for _, dbName := range databaseNames {
		if dbName == "kratos" {
			kratosDatabaseExists = true
			break
		}
	}

	// 如果不存在，创建数据库
	if !kratosDatabaseExists {
		// 创建 "kratos" 数据库
		err := client.Database("kratos").CreateCollection(context.TODO(), "dummy")
		if err != nil {
			return nil, fmt.Errorf("failed to create kratos database: %w", err)
		}
	}

	// 选择数据库 my_db
	db = client.Database("kratos")
	return db, nil
}
