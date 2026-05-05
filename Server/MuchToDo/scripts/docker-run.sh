#!/bin/bash
set -euo pipefail

ACTION="${1:-up}"

echo "======================================"
echo "  Docker Compose: $ACTION"
echo "======================================"

case "$ACTION" in
  up)
    docker compose up --build -d
    echo ""
    echo "✅ Services started. Waiting for health checks..."
    sleep 5
    docker compose ps
    echo ""
    echo "API:     http://localhost:8080"
    echo "Health:  http://localhost:8080/health"
    echo "MongoDB: localhost:27017"
    ;;
  down)
    docker compose down
    echo "✅ Services stopped."
    ;;
  logs)
    docker compose logs -f
    ;;
  restart)
    docker compose restart
    ;;
  *)
    echo "Usage: $0 [up|down|logs|restart]"
    exit 1
    ;;
esac
