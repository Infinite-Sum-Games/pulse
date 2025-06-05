# Pulse - Monolithic Backend for ACM's Season of Code

This repository contains the monolithic Golang backend to ACM's
Season of Code platform. 

### Setup Steps

1. Clone the repository
```bash
git clone https://github.com/Infinite-Sum-Games/pulse
# or
gh repo clone Infinite-Sum-Games/pulse
```
2. Setup a PostgreSQL database. You can either use [Neon](https://neon.tech) 
and grab the connection string from there or use `docker compose`. If you are 
using compose, then run the following command:
```bash
# For newer versions with Docker Desktop available
docker compose up -d
# For older versions with docker-cli available
docker-compose up -d
```
This should start a container named `pulse-postgres`. Your database connection 
string is the default string provided in `.env.example`

2. Fill out the `environment variables` and rename the file `.env.example` as 
`.env`.

3. Run the database migrations as follows:
```bash
# Download the tool from https://github.com/pressly/goose
# 1. For creating tables
make up
# 2. For adding seed data
make seed
# 3. For cleaning and deleting all tables
make down
```
4. Generate all the database helper functions by running:
```bash
# Download the tool from https://github.com/sqlc-dev/sqlc
sqlc generate
```
5. For development you can get live-reloading features by using:
```bash
# Download the tool from https://github.com/air-verse/air
air
```
6. For building the project and running use:
```bash
# Install make in your environment
make run
# Alternatively for running without building
go run main.go
# For getting a build and then running
go build -o bin/pulse
./bin/pulse
```

> The development of this repository has taken place in a linux environment. 
It would be easier if a Linux (or) Unix (or) WSL environment is used. Otherwise
the setup instructions remain the same, but need to be tailored to Windows.

### Testing Instructions
1. For testing the API, download [Bruno - API Client](https://www.usebruno.com/)
2. Open the `bruno/` folder in this repository using Bruno.
3. Run individual requests.
4. For checking LIVE updates, check `Makefile` and populate Valkey / Redis 
accordingly.
5. A sample event to test the LIVE update are as follows:
```bash
# Use the Redis-CLI or Redis-Insight on localhost:5540 to add the following
# Add a few sample events to the live-update-stream

# Sample 1
redis-cli XADD live-update-stream * data '{
  "github_username": "IAmRiteshKoushik",
  "message": "IAmRiteshKoushik claimed a bounty worth 500 points!", 
  "event_type": "Bounty",
  "time": 1703764533000
}'

# Sample 2
redis-cli XADD live-update-stream * data '{
  "github_username": "vijay-sb",
  "message": "vijay-sb claimed a bounty worth 700 points!", 
  "event_type": "Bounty",
  "time": 1703764533000
}'
```



### Authors
This application has been built and tested by [Ritesh Koushik](https://github.com/IAmRiteshKoushik).
