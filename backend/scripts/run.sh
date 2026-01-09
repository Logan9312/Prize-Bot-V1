docker build -t discord-auction-bot .
docker run --env-file .env -p 8080:8080 discord-auction-bot