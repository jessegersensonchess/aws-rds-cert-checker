# AWS RDS Certificate Retriever

This application retrieves the Certificate Authority (CA) for each database instance in AWS using the RDS service.

## Prerequisites

- Docker installed on your machine.

## Setup & Run

### 1. Clone the repository

```bash
git clone https://github.com/jessegersenson/aws-rds-cert-checker
cd aws-rds-cert-checker
```

### 2. Build the Docker Image
```
docker build -t aws-rds-cert-checker .
```

### 3. Run the Docker Container
```
docker run --rm -v $HOME/.aws:/root/.aws aws-rds-cert-checker:latest
```
By default, the app looks for databases in the us-west-1 region with the default AWS profile. To specify different regions or profiles, modify the Go program accordingly and rebuild the Docker image.



