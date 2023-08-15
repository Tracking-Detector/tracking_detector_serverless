#!/bin/bash

RED=$(tput setaf 1)
GREEN=$(tput setaf 2)
YELLOW=$(tput setaf 3)
NC=$(tput sgr0) 

generate_random_string() {
    openssl rand -base64 48 | tr -dc 'a-zA-Z0-9' | head -c $1
}

prompt_or_generate() {
  local description=$1
  local length=$2

  echo -e "${YELLOW}$description:${NC}" >&2
  echo -e "1. ${GREEN}Generate a random value (recommended)${NC}" >&2
  echo -e "2. ${RED}Input your own value${NC}" >&2

  read -p "Choose option (1/2) > " choice >&2
  if [[ $choice -eq 1 ]]; then
    generate_random_string $length
  else
    read -p "Enter $description: " value >&2
    echo $value
  fi
}

echo -e "${GREEN}Welcome to the Tracking Detector Gym setup wizard!${NC}\n"

admin_api_key=$(prompt_or_generate "Admin API Key" 32)
admin_password=$(prompt_or_generate "Admin Password" 32)
minio_private_key=$(prompt_or_generate "MinIO Private Key" 32)

# Write to .env file
cat > .env <<EOL
# mongo
MONGO_URI=mongodb://db:27017/tracking-detector
USER_COLLECTION=users
REQUEST_COLLECTION=requests
TRAINING_RUNS_COLLECTION=training-runs

# minio
MINIO_URI=minio:9000
MINIO_ACCESS_KEY=adminadmin
MINIO_PRIVATE_KEY=$minio_private_key
EXPORT_BUCKET_NAME=exports
MODEL_BUCKET_NAME=models

# authentication
ADMIN_API_KEY=$admin_api_key
ADMIN_USERNAME=admin
ADMIN_PASSWORD=$admin_password
EOL

echo -e "\n${GREEN}Setup completed and .env file generated!${NC}\n"
echo "Generated/Entered values:"
echo -e "${YELLOW}Admin API Key:${NC} $admin_api_key"
echo -e "${YELLOW}Admin Password:${NC} $admin_password"
echo -e "${YELLOW}MinIO Private Key:${NC} $minio_private_key"

read -p "${GREEN}Do you want to start Docker Compose (build and up)? (y/n):${NC} " start_docker

if [[ $start_docker == 'y' || $start_docker == 'Y' ]]; then
    echo -e "${GREEN}Starting Docker Compose...${NC}"
    docker-compose build
    docker-compose up -d
    echo -e "${GREEN}Docker Compose has been started in daemon mode.${NC}"
else
    echo -e "${RED}Exiting without starting Docker Compose.${NC}"
fi