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

func EnvLogStashHost() string {
	return os.Getenv("LOGSTASH_HOST")
}

func EnvLogStashPort() string {
	return os.Getenv("LOGSTASH_PORT")
}
