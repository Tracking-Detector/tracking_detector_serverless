package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

func ConnectMinio() *minio.Client {
	minioClient, err := minio.New(EnvMinIoURI(), &minio.Options{
		Creds:  credentials.NewStaticV4(EnvMinIoAccessKey(), EnvMinIoPrivateKey(), ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connected to MinIO")
	return minioClient
}

var DB *mongo.Client = ConnectDB()

var MINIO *minio.Client = ConnectMinio()

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("tracking-detector").Collection(collectionName)
	return collection
}

func VerifyBucketExists(ctx context.Context, client *minio.Client, bucketName string) {
	if exists, err := client.BucketExists(ctx, bucketName); err != nil {
		log.Fatal(err)
	} else if exists {
	} else {
		if makeBucketError := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "eu-central-1"}); makeBucketError != nil {
			log.Fatal(makeBucketError)
		} else {
			if setVersioningError := client.SetBucketVersioning(ctx, bucketName, minio.BucketVersioningConfiguration{
				Status: "Enabled",
			}); setVersioningError != nil {
				log.Fatal(setVersioningError)
			}
		}
	}
}
