output "ecs_exec_role_policy_id" {
  description = "The ECS service role policy ID, in the form of `role_name:role_policy_name`"
  value       = module.ecs_alb_service_task.ecs_exec_role_policy_id
}

output "ecs_exec_role_policy_name" {
  description = "ECS service role name"
  value       = module.ecs_alb_service_task.ecs_exec_role_policy_name
}

output "service_name" {
  description = "ECS Service name"
  value       = module.ecs_alb_service_task.service_name
}

output "service_role_arn" {
  description = "ECS Service role ARN"
  value       = module.ecs_alb_service_task.service_role_arn
}

output "task_exec_role_name" {
  description = "ECS Task role name"
  value       = module.ecs_alb_service_task.task_exec_role_name
}

output "task_exec_role_arn" {
  description = "ECS Task exec role ARN"
  value       = module.ecs_alb_service_task.task_exec_role_arn
}

output "task_role_name" {
  description = "ECS Task role name"
  value       = module.ecs_alb_service_task.task_role_name
}

output "task_role_arn" {
  description = "ECS Task role ARN"
  value       = module.ecs_alb_service_task.task_role_arn
}

output "task_role_id" {
  description = "ECS Task role id"
  value       = module.ecs_alb_service_task.task_role_id
}

output "service_security_group_id" {
  description = "Security Group ID of the ECS task"
  value       = module.ecs_alb_service_task.service_security_group_id
}

output "task_definition_family" {
  description = "ECS task definition family"
  value       = module.ecs_alb_service_task.task_definition_family
}

output "task_definition_revision" {
  description = "ECS task definition revision"
  value       = module.ecs_alb_service_task.task_definition_revision
}

output "bucket_domain_name" {
  value       = var.s3_enabled ? module.s3_bucket.bucket_domain_name : ""
  description = "FQDN of bucket"
}

output "bucket_regional_domain_name" {
  value       = var.s3_enabled ? module.s3_bucket.bucket_regional_domain_name : ""
  description = "The bucket region-specific domain name"
}

output "bucket_id" {
  value       = var.s3_enabled ? module.s3_bucket.bucket_id : ""
  description = "Bucket Name (aka ID)"
}

output "bucket_arn" {
  value       = var.s3_enabled ? module.s3_bucket.bucket_arn : ""
  description = "Bucket ARN"
}

output "s3_user_arn" {
  value       = var.s3_user_enabled ? module.s3_bucket.user_arn : ""
  description = "User ARN"
}

output "s3_user_name" {
  value       = var.s3_user_enabled ? module.s3_bucket.user_name : ""
  description = "User name"
}

output "s3_user_unique_id" {
  value       = var.s3_user_enabled ? module.s3_bucket.user_unique_id : ""
  description = "User unique ID"
}

output "rds_instance_id" {
  value       = module.rds_instance.instance_id
  description = "ID of the instance"
}

output "rds_instance_address" {
  value       = module.rds_instance.instance_address
  description = "Address of the instance"
}

output "rds_instance_endpoint" {
  value       = module.rds_instance.instance_endpoint
  description = "DNS Endpoint of the instance"
}

output "rds_subnet_group_id" {
  value       = module.rds_instance.subnet_group_id
  description = "ID of the created Subnet Group"
}

output "rds_security_group_id" {
  value       = module.rds_instance.security_group_id
  description = "ID of the Security Group"
}

output "rds_parameter_group_id" {
  value       = module.rds_instance.parameter_group_id
  description = "ID of the Parameter Group"
}

output "rds_option_group_id" {
  value       = module.rds_instance.option_group_id
  description = "ID of the Option Group"
}

output "rds_hostname" {
  value       = module.rds_instance.hostname
  description = "DNS host name of the instance"
}
