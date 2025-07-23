# Local Development Guide

This guide explains how to run all backend services locally for development.

## Prerequisites

Before running the services locally, ensure you have the following dependencies running:

### Required External Services

1. **MySQL Database**
   - Host: `localhost:53306`
   - Username: `root`
   - Password: `root`
   - Database: `golden_pay_db`

2. **Kafka**
   - Host: `localhost:59092`
   - Version: 3.6.0

### Setting up Dependencies

You can use Docker to set up the required dependencies:

```bash
# Start MySQL
docker run -d \
  --name goldenpay-mysql \
  -p 53306:3306 \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=golden_pay_db \
  mysql:8.0

# Start Kafka (requires Zookeeper)
docker run -d \
  --name goldenpay-zookeeper \
  -p 2181:2181 \
  confluentinc/cp-zookeeper:latest \
  bash -c "echo 'clientPort=2181' > /etc/kafka/zookeeper.properties && \
           echo 'dataDir=/var/lib/zookeeper' >> /etc/kafka/zookeeper.properties && \
           /etc/confluent/docker/run"

docker run -d \
  --name goldenpay-kafka \
  -p 59092:59092 \
  -e KAFKA_ZOOKEEPER_CONNECT=localhost:2181 \
  -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:59092 \
  -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 \
  --network=host \
  confluentinc/cp-kafka:latest
```

## Services Architecture

The backend consists of 5 microservices:

| Service | Type | Port | Dependencies |
|---------|------|------|-------------|
| **user** | gRPC | 5001 | MySQL, Kafka |
| **wallet** | gRPC | 5002 | MySQL, Kafka |
| **chat** | gRPC | 5003 | MySQL, Kafka |
| **http** | HTTP API | 5000 | user, wallet services |
| **event-handler** | Kafka Consumer | - | user, wallet services, Kafka |

## Running Services Locally

### Using the Local Run Script

The easiest way to run all services is using the provided script:

```bash
# Start all services
make run/local
# OR
./script/run_local.sh start

# Check status
make status/local
# OR
./script/run_local.sh status

# View logs
make logs/local
# OR
./script/run_local.sh logs

# Stop all services
make stop/local
# OR
./script/run_local.sh stop

# Restart all services
make restart/local
# OR
./script/run_local.sh restart

# Clean logs and pid files
make clean/local
# OR
./script/run_local.sh clean
```

### Manual Service Management

If you need to manage services individually:

```bash
# Build all services
make build/all

# Build individual service
make build/user
make build/wallet
make build/chat
make build/http
make build/event-handler

# Run individual service (manual)
cd internal/service/user/config && ../../../bin/user/exc
```

## Service Startup Order

The script automatically handles dependency order:

1. **Independent Services** (started first):
   - `user` service (port 5001)
   - `wallet` service (port 5002)
   - `chat` service (port 5003)

2. **Dependent Services** (started after 5-second delay):
   - `http` service (port 5000) - depends on user, wallet
   - `event-handler` - depends on user, wallet

## Accessing Services

Once all services are running:

- **HTTP API**: http://localhost:5000
  - REST endpoints for frontend integration
  
- **gRPC Services**:
  - User Service: `localhost:5001`
  - Wallet Service: `localhost:5002`
  - Chat Service: `localhost:5003`

## Logs and Debugging

### Log Files

All service logs are stored in `./logs/` directory:
- `./logs/user.log`
- `./logs/wallet.log`
- `./logs/chat.log`
- `./logs/http.log`
- `./logs/event-handler.log`

### Viewing Logs

```bash
# View all logs
make logs/local

# View individual service logs
tail -f logs/user.log
tail -f logs/http.log

# View logs in real-time for all services
tail -f logs/*.log
```

### Process Management

Service PIDs are tracked in `./local_services.pid`. The script handles:
- Graceful shutdown (SIGTERM)
- Force kill if needed (SIGKILL)
- Process monitoring
- Automatic restart detection

## Common Issues

### Port Already in Use
If you get "port already in use" errors:
```bash
# Check what's using the port
lsof -i :5001
lsof -i :5000

# Stop all local services
make stop/local
```

### Database Connection Issues
1. Ensure MySQL is running on `localhost:53306`
2. Check credentials: `root/root`
3. Ensure database `golden_pay_db` exists

### Kafka Connection Issues
1. Ensure Kafka is running on `localhost:59092`
2. Check Zookeeper is running on `localhost:2181`
3. Verify Kafka topics are created

### Service Won't Start
1. Check logs: `make logs/local`
2. Verify dependencies are running
3. Ensure binaries are built: `make build/all`
4. Check config files in each service's `config/` directory

## Development Workflow

1. **Start Dependencies**: MySQL, Kafka
2. **Start Services**: `make run/local`
3. **Check Status**: `make status/local`
4. **Make Changes**: Edit code
5. **Restart**: `make restart/local` (rebuilds automatically)
6. **Debug**: Check logs with `make logs/local`
7. **Stop**: `make stop/local` when done

## Configuration

Each service reads its configuration from:
- `./internal/service/{service}/config/config.yaml`

Key configuration sections:
- Database connection
- Service addresses and ports
- Kafka settings
- JWT settings (for user service)

For production deployment, use the corresponding `container_config.yaml` files. 