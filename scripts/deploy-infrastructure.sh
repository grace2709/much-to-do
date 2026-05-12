#!/bin/bash
set -e

echo "🏗️ Starting Infrastructure Deployment"

# Validate Terraform
cd terraform
terraform fmt -check
terraform validate

# Plan and Apply
terraform plan -out=tfplan
terraform apply tfplan

# Export outputs for application deployment
terraform output -json > ../infrastructure-outputs.json

echo "✅ Infrastructure deployed successfully"
