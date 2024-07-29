# articles-system

## Description
articles-system is a Golang web-API application that allows users to post articles and retrieve articles via HTTP.

## Needed Stacks
1. Golang
2. MySQL
3. Redis (standalone)

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/Bainandhika/articles-system
   cd articles-system
2. Install dependencies:
   ```bash
   go mod tidy
3. Prepare the .env file:
   Create a .env file in the project root directory with the following content:
   ```env
   APP_HOST=localhost
   APP_PORT=8080
   APP_LOG_PATH=logs
   
   DATABASE_HOST=127.0.0.1
   DATABASE_PORT=3306
   DATABASE_USERNAME=root
   DATABASE_PASSWORD=
   DATABASE_NAME=articles_system
   
   REDIS_HOST=127.0.0.1
   REDIS_PORT=6379
   REDIS_USERNAME=
   REDIS_PASSWORD=
4. Build the executable:
   ```bash
   go build
5. Make sure MySQL and Redis are running and the database is already created. Then run the application:
   ```bash
   ./main

## Postman Collection
You can find the collection at assets directory

## Author
M. Bainandhika Baghaskara Putra - https://github.com/Bainandhika
  
