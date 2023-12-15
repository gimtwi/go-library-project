# Go Application with PostgreSQL Database

This repository contains a Go application that uses a PostgreSQL database running in a Docker container.

## Prerequisites

Before you can run this application, make sure you have the following prerequisites installed on your system:

- Go (https://golang.org/)
- Docker (https://www.docker.com/)

## Setup

1. Clone this repository to your local machine:

   ```bash
   git clone https://github.com/yourusername/your-app.git
   cd your-app

2. Run following command to start a postgres database:

   ```bash
   docker run --name go-library -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=password -e POSTGRES_DB=library-project -p 5432:5432 -d postgres

3. Same for the test database:

   ```bash
   docker run --name go-library-test -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=password -e POSTGRES_DB=library-project-test -p 5433:5432 -d postgres

