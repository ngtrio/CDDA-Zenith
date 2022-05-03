docker-compose down
cp -r /ngtrio.me-ssl-bundle ./nginx/
docker-compose up --build -d