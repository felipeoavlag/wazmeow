#!/bin/bash

# WazMeow Docker Management Script
# Usage: ./scripts/docker.sh [command]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker first."
        exit 1
    fi
}

# Function to check if docker-compose is available
check_docker_compose() {
    if ! command -v docker-compose > /dev/null 2>&1; then
        print_error "docker-compose is not installed or not in PATH."
        exit 1
    fi
}

# Function to start services
start_services() {
    print_info "Starting WazMeow services..."
    
    # Remove any existing containers
    docker-compose down --remove-orphans 2>/dev/null || true
    
    # Build and start services
    docker-compose up -d --build
    
    print_info "Waiting for services to be ready..."
    
    # Wait for PostgreSQL to be ready
    print_info "Waiting for PostgreSQL to be ready..."
    timeout=60
    counter=0
    
    while [ $counter -lt $timeout ]; do
        if docker-compose exec -T postgres pg_isready -U wazmeow -d wazmeow > /dev/null 2>&1; then
            print_success "PostgreSQL is ready!"
            break
        fi
        
        if [ $counter -eq $((timeout - 1)) ]; then
            print_error "PostgreSQL failed to start within $timeout seconds"
            print_info "Checking PostgreSQL logs..."
            docker-compose logs postgres
            exit 1
        fi
        
        sleep 1
        counter=$((counter + 1))
        echo -n "."
    done
    
    echo ""
    print_success "All services are running!"
    print_info "WazMeow API is available at: http://localhost:8080"
    print_info "PostgreSQL is available at: localhost:5432"
}

# Function to stop services
stop_services() {
    print_info "Stopping WazMeow services..."
    docker-compose down
    print_success "Services stopped!"
}

# Function to restart services
restart_services() {
    print_info "Restarting WazMeow services..."
    stop_services
    start_services
}

# Function to show logs
show_logs() {
    service=${1:-""}
    if [ -z "$service" ]; then
        print_info "Showing logs for all services..."
        docker-compose logs -f
    else
        print_info "Showing logs for service: $service"
        docker-compose logs -f "$service"
    fi
}

# Function to show status
show_status() {
    print_info "Service status:"
    docker-compose ps
    
    echo ""
    print_info "Health checks:"
    
    # Check PostgreSQL
    if docker-compose exec -T postgres pg_isready -U wazmeow -d wazmeow > /dev/null 2>&1; then
        print_success "PostgreSQL: Healthy"
    else
        print_error "PostgreSQL: Unhealthy"
    fi
    
    # Check WazMeow API
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        print_success "WazMeow API: Healthy"
    else
        print_warning "WazMeow API: Not responding (may still be starting)"
    fi
}

# Function to clean up
cleanup() {
    print_info "Cleaning up Docker resources..."
    docker-compose down --volumes --remove-orphans
    docker system prune -f
    print_success "Cleanup completed!"
}

# Function to enter database shell
db_shell() {
    print_info "Connecting to PostgreSQL database..."
    docker-compose exec postgres psql -U wazmeow -d wazmeow
}

# Function to show help
show_help() {
    echo "WazMeow Docker Management Script"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  start     Start all services"
    echo "  stop      Stop all services"
    echo "  restart   Restart all services"
    echo "  status    Show service status"
    echo "  logs      Show logs for all services"
    echo "  logs <service>  Show logs for specific service (postgres, wazmeow)"
    echo "  db        Connect to PostgreSQL database"
    echo "  cleanup   Stop services and clean up Docker resources"
    echo "  help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 start"
    echo "  $0 logs postgres"
    echo "  $0 status"
}

# Main script logic
main() {
    # Check prerequisites
    check_docker
    check_docker_compose
    
    # Parse command
    command=${1:-"help"}
    
    case $command in
        "start")
            start_services
            ;;
        "stop")
            stop_services
            ;;
        "restart")
            restart_services
            ;;
        "status")
            show_status
            ;;
        "logs")
            show_logs "$2"
            ;;
        "db")
            db_shell
            ;;
        "cleanup")
            cleanup
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            print_error "Unknown command: $command"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
