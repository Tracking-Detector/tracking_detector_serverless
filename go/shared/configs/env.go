package configs

import (
	"os"
)

func EnvMongoURI() string {
	return os.Getenv("MONGO_URI")
}

func EnvMinIoURI() string {
	return os.Getenv("MINIO_URI")
}

func EnvMinIoAccessKey() string {
	return os.Getenv("MINIO_ACCESS_KEY")
}

func EnvMinIoPrivateKey() string {
	return os.Getenv("MINIO_PRIVATE_KEY")
}

func EnvExportBucketName() string {
	return os.Getenv("EXPORT_BUCKET_NAME")
}

func EnvModelBucketName() string {
	return os.Getenv("MODEL_BUCKET_NAME")
}
