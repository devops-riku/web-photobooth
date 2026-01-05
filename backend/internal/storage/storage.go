package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	appconfig "web-photobooth/backend/internal/config"
)

func InitDB(cfg *appconfig.Config) (*gorm.DB, error) {
	// 1. Create pgxpool configuration
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// 2. Tune pool settings for performance
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	// 3. Create the pool
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	// 4. Ping to verify connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	// 5. Wrap pgxpool as an sql.DB to use with GORM
	sqlDB := stdlib.OpenDBFromPool(pool)

	// 6. Initialize GORM with the existing connection
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitS3(cfg *appconfig.Config) (*s3.Client, error) {
	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               cfg.DOEndpoint,
			SigningRegion:     cfg.DORegion,
			HostnameImmutable: true,
		}, nil
	})

	s3Cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(cfg.DORegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.DOKey, cfg.DOSecret, "")),
		awsconfig.WithEndpointResolverWithOptions(resolver),
	)
	if err != nil {
		return nil, err
	}
	return s3.NewFromConfig(s3Cfg), nil
}
