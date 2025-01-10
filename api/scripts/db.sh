docker stop photobox-db; \
docker rm photobox-db; \
docker run \
  -d \
  --name photobox-db \
  -e "POSTGRES_USER=admin" \
  -e "POSTGRES_PASSWORD=password" \
  -p 5432:5432 \
  -v $(pwd)/sql:/docker-entrypoint-initdb.d \
  postgres:17.2