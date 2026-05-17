output "cloudfront_domain" {
  value = module.storage.cloudfront_domain
}

output "alb_dns_name" {
  value = module.compute.alb_dns_name
}

output "redis_endpoint" {
  value     = module.storage.redis_endpoint
  sensitive = true
}

output "s3_bucket_name" {
  value = module.storage.s3_bucket_name
