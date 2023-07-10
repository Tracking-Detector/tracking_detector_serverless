# tracking_detector_serverless


# Setup

It is important that you change the usernames and passwords in the .env file especially:
```sh
# mongo
MONGO_URI=mongodb://db:27017/tracking-detector
USER_COLLECTION=users
REQUEST_COLLECTION=requests
TRAINING_RUNS_COLLECTION=training-runs
# minio
MINIO_URI=minio:9000
# Set this
MINIO_ACCESS_KEY=
# Set this
MINIO_PRIVATE_KEY=
EXPORT_BUCKET_NAME=exports
MODEL_BUCKET_NAME=models

# Set this
ADMIN_API_KEY=
```
Also in the dir /infra/api-gateway create a .htpasswd file with with a username and a password
```sh
# Maybe install
sudo apt-get install apache2-utils
# Generate file
htpasswd -c -B -b infra/api-gateway/.htpasswd <username> <password>
```
After that run: 
```sh
docker compose build
docker compose up
```