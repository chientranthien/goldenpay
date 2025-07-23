#!/usr/bin/env bash

# Get project root directory (where this script is located relative to)
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Source common functions using absolute path
. "$PROJECT_ROOT/script/common.sh"

# Configuration
SERVICES=("user" "wallet" "chat" "http" "event-handler")
INDEPENDENT_SERVICES=("user" "wallet" "chat")
DEPENDENT_SERVICES=("http" "event-handler")
LOG_DIR="$PROJECT_ROOT/logs"
PID_FILE="$PROJECT_ROOT/local_services.pid"

# Colors for logging
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

usage() {
  cat <<EOF
Run all services locally for development.

Usage:
  ./script/run_local.sh [COMMAND]

Commands:
  start     Start all services (default)
  stop      Stop all running services
  restart   Restart all services
  status    Show status of all services
  logs      Show logs for all services
  clean     Clean logs and pid files

Dependencies:
  - MySQL running on localhost:53306
  - Kafka running on localhost:59092

Services will be started in dependency order:
  1. Independent: user, wallet, chat
  2. Dependent: http, event-handler
EOF
}

# Logging functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_service() {
    echo -e "${BLUE}[$1]${NC} $2"
}

# Initialize logging directory
init_logging() {
    mkdir -p "$LOG_DIR"
    if [ ! -d "$LOG_DIR" ]; then
        log_error "Failed to create log directory: $LOG_DIR"
        exit 1
    fi
}

# Check if service binary exists
check_binary() {
    local service=$1
    local binary="$PROJECT_ROOT/bin/$service/exc"
    
    if [ ! -f "$binary" ]; then
        log_warn "Binary for $service not found. Building..."
        cd "$PROJECT_ROOT"
        make build/$service
        if [ $? -ne 0 ]; then
            log_error "Failed to build $service"
            return 1
        fi
    fi
    return 0
}

# Start a single service
start_service() {
    local service=$1
    local binary="$PROJECT_ROOT/bin/$service/exc"
    local log_file="$LOG_DIR/$service.log"
    local config_dir="$PROJECT_ROOT/internal/service/$service/config"
    
    # Check if service is already running
    if is_service_running "$service"; then
        log_warn "$service is already running"
        return 0
    fi
    
    # Check binary exists
    if ! check_binary "$service"; then
        return 1
    fi
    
    # Change to service config directory and start service
    log_info "Starting $service service..."
    cd "$config_dir"
    
    # Start service in background and redirect output to log file
    nohup "$binary" > "$log_file" 2>&1 &
    local pid=$!
    
    # Save PID
    echo "$service:$pid" >> "$PID_FILE"
    
    # Wait a moment to check if service started successfully
    sleep 2
    if kill -0 $pid 2>/dev/null; then
        log_service "$service" "Started successfully (PID: $pid)"
        return 0
    else
        log_error "Failed to start $service. Check logs: $log_file"
        return 1
    fi
}

# Stop a single service
stop_service() {
    local service=$1
    local pid=$(get_service_pid "$service")
    
    if [ -z "$pid" ]; then
        log_warn "$service is not running"
        return 0
    fi
    
    log_info "Stopping $service service (PID: $pid)..."
    
    # Try graceful shutdown first
    if kill -TERM $pid 2>/dev/null; then
        # Wait up to 10 seconds for graceful shutdown
        local count=0
        while [ $count -lt 10 ] && kill -0 $pid 2>/dev/null; do
            sleep 1
            count=$((count + 1))
        done
        
        # Force kill if still running
        if kill -0 $pid 2>/dev/null; then
            log_warn "Force killing $service..."
            kill -KILL $pid 2>/dev/null
        fi
        
        log_service "$service" "Stopped"
    else
        log_warn "$service process not found or already stopped"
    fi
    
    # Remove from PID file
    if [ -f "$PID_FILE" ]; then
        grep -v "^$service:" "$PID_FILE" > "$PID_FILE.tmp"
        mv "$PID_FILE.tmp" "$PID_FILE"
    fi
}

# Get service PID
get_service_pid() {
    local service=$1
    if [ -f "$PID_FILE" ]; then
        grep "^$service:" "$PID_FILE" | cut -d: -f2
    fi
}

# Check if service is running
is_service_running() {
    local service=$1
    local pid=$(get_service_pid "$service")
    
    if [ -n "$pid" ] && kill -0 $pid 2>/dev/null; then
        return 0
    else
        return 1
    fi
}

# Start all services
start_all() {
    log_info "Building all services..."
    cd "$PROJECT_ROOT"
    make build/all
    if [ $? -ne 0 ]; then
        log_error "Build failed. Aborting."
        exit 1
    fi
    
    init_logging
    
    # Clear PID file
    > "$PID_FILE"
    
    log_info "Starting independent services first..."
    for service in "${INDEPENDENT_SERVICES[@]}"; do
        start_service "$service"
        if [ $? -ne 0 ]; then
            log_error "Failed to start $service. Check dependencies."
            exit 1
        fi
    done
    
    # Wait a bit for independent services to be ready
    log_info "Waiting for independent services to be ready..."
    sleep 5
    
    log_info "Starting dependent services..."
    for service in "${DEPENDENT_SERVICES[@]}"; do
        start_service "$service"
        if [ $? -ne 0 ]; then
            log_error "Failed to start $service"
            exit 1
        fi
    done
    
    log_info "All services started successfully!"
    echo ""
    show_status
}

# Stop all services
stop_all() {
    if [ ! -f "$PID_FILE" ]; then
        log_info "No services are running"
        return 0
    fi
    
    log_info "Stopping all services..."
    
    # Stop dependent services first
    for service in "${DEPENDENT_SERVICES[@]}"; do
        stop_service "$service"
    done
    
    # Then stop independent services
    for service in "${INDEPENDENT_SERVICES[@]}"; do
        stop_service "$service"
    done
    
    # Clean up PID file
    rm -f "$PID_FILE"
    
    log_info "All services stopped"
}

# Show status of all services
show_status() {
    echo "Service Status:"
    echo "==============="
    
    for service in "${SERVICES[@]}"; do
        if is_service_running "$service"; then
            local pid=$(get_service_pid "$service")
            echo -e "${GREEN}✓${NC} $service (PID: $pid)"
        else
            echo -e "${RED}✗${NC} $service"
        fi
    done
    
    echo ""
    echo "Service URLs:"
    echo "============="
    if is_service_running "http"; then
        echo -e "${BLUE}HTTP API:${NC} http://localhost:5000"
    fi
    if is_service_running "user"; then
        echo -e "${BLUE}User gRPC:${NC} localhost:5001"
    fi
    if is_service_running "wallet"; then
        echo -e "${BLUE}Wallet gRPC:${NC} localhost:5002"
    fi
    if is_service_running "chat"; then
        echo -e "${BLUE}Chat gRPC:${NC} localhost:5003"
    fi
}

# Show logs for all services
show_logs() {
    if [ ! -d "$LOG_DIR" ]; then
        log_error "No logs directory found"
        return 1
    fi
    
    for service in "${SERVICES[@]}"; do
        local log_file="$LOG_DIR/$service.log"
        if [ -f "$log_file" ]; then
            echo -e "${BLUE}=== $service logs ===${NC}"
            tail -n 20 "$log_file"
            echo ""
        fi
    done
}

# Clean logs and pid files
clean() {
    log_info "Cleaning logs and pid files..."
    rm -rf "$LOG_DIR"
    rm -f "$PID_FILE"
    log_info "Cleanup completed"
}

# Signal handler for graceful shutdown
cleanup_on_exit() {
    echo ""
    log_info "Received interrupt signal. Stopping all services..."
    stop_all
    exit 0
}

# Set signal handler
trap cleanup_on_exit SIGINT SIGTERM

# Main command processing
case "${1:-start}" in
    start)
        start_all
        ;;
    stop)
        stop_all
        ;;
    restart)
        stop_all
        sleep 2
        start_all
        ;;
    status)
        show_status
        ;;
    logs)
        show_logs
        ;;
    clean)
        clean
        ;;
    -h|--help)
        usage
        ;;
    *)
        log_error "Unknown command: $1"
        usage
        exit 1
        ;;
esac 