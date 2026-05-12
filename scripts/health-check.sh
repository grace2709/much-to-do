#!/bin/bash

echo "🔍 Running Health Checks..."

# Check ALB
ALB_URL="${ALB_DNS_NAME}"
if curl -f -s "${ALB_URL}/health" > /dev/null; then
    echo "✅ ALB health check passed"
else
    echo "❌ ALB health check failed"
    exit 1
fi

# Check Frontend
CF_URL="${CLOUDFRONT_DOMAIN}"
if curl -f -s "https://${CF_URL}" > /dev/null; then
    echo "✅ Frontend accessible"
else
    echo "❌ Frontend check failed"
    exit 1
fi

echo "✅ All health checks passed"
