provider "aws" {
  region = var.region
}

# Create label
module "label" {
  source     = "git::https://github.com/cloudposse/terraform-null-label.git?ref=tags/0.16.0"
  namespace  = var.namespace
  stage      = var.stage
  name       = var.name
  attributes = var.attributes
  delimiter  = var.delimiter

  tags = var.tags
}

# Create Target Group
resource "aws_lb_target_group" "default" {
  name                 = module.label.id
  port                 = var.target_group_port
  protocol             = var.target_group_protocol
  target_type          = var.target_group_target_type
  vpc_id               = var.vpc_id
  deregistration_delay = var.deregistration_delay

  health_check {
    protocol            = var.target_group_protocol
    path                = var.health_check_path
    timeout             = var.health_check_timeout
    healthy_threshold   = var.health_check_healthy_threshold
    unhealthy_threshold = var.health_check_unhealthy_threshold
    interval            = var.health_check_interval
    matcher             = var.health_check_matcher
  }

  tags = var.tags
}

# Add ALB listener rule
resource "aws_lb_listener_rule" "static" {
  listener_arn = var.load_balancer_listener_arn

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.default.arn
  }

  condition {
    path_pattern {
      values = var.load_balancer_listener_paths
    }
  }
}

module "s3_bucket" {
  source                 = "git::https://github.com/cloudposse/terraform-aws-s3-bucket.git?ref=tags/0.16.0"
  acl                    = "private"
  enabled                = var.s3_enabled
  user_enabled           = var.s3_user_enabled
  versioning_enabled     = false
  allowed_bucket_actions = ["s3:*"]
  name                   = var.name
  stage                  = var.stage
  namespace              = var.namespace
}

resource "aws_cloudwatch_log_group" "default" {
  name = module.label.id

  tags = var.tags
}

resource "aws_ecr_repository" "default" {
  name                 = module.label.name
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = false
  }
}

# Create Container
module "container_definition" {
  source                       = "git::https://github.com/cloudposse/terraform-aws-ecs-container-definition.git?ref=tags/0.21.0"
  container_name               = module.label.id
  container_image              = format("%s:%s", module.aws_ecr_repository.default.repository_url, module.label.stage)
  container_memory             = var.container_memory
  container_memory_reservation = var.container_memory_reservation
  container_cpu                = var.container_cpu
  essential                    = var.container_essential
  readonly_root_filesystem     = var.container_readonly_root_filesystem
  environment                  = var.container_environment
  secrets                      = var.container_secrets
  port_mappings                = var.container_port_mappings
  log_configuration = {
    logDriver = "awslogs"
    options = {
      "awslogs-region"        = var.region
      "awslogs-group"         = module.aws_cloudwatch_log_group.default.name
      "awslogs-stream-prefix" = "ecs"
    }
    secretOptions = null
  }
}

# Create Service Task
module "ecs_alb_service_task" {
  source                             = "git::https://github.com/cloudposse/terraform-aws-ecs-alb-service-task.git?ref=tags/0.31.0"
  namespace                          = var.namespace
  stage                              = var.stage
  name                               = var.name
  attributes                         = var.attributes
  delimiter                          = var.delimiter
  alb_security_group                 = var.vpc_default_security_group_id
  container_definition_json          = module.container_definition.json
  ecs_cluster_arn                    = var.ecs_cluster_arn
  launch_type                        = var.ecs_launch_type
  vpc_id                             = var.vpc_id
  security_group_ids                 = [var.vpc_default_security_group_id]
  subnet_ids                         = var.private_subnet_ids
  tags                               = var.tags
  ignore_changes_task_definition     = var.ignore_changes_task_definition
  network_mode                       = var.network_mode
  assign_public_ip                   = var.assign_public_ip
  propagate_tags                     = var.propagate_tags
  deployment_minimum_healthy_percent = var.deployment_minimum_healthy_percent
  deployment_maximum_percent         = var.deployment_maximum_percent
  deployment_controller_type         = var.deployment_controller_type
  desired_count                      = var.desired_count
  task_memory                        = var.task_memory
  task_cpu                           = var.task_cpu
  ecs_load_balancers = [
    {
      container_name   = module.label.id
      container_port   = var.container_port
      elb_name         = ""
      target_group_arn = aws_lb_target_group.default.arn
    }
  ]
}

# Attach task roles
resource "aws_iam_policy" "task_role_policy" {
  name   = "${module.label.id}-task"
  policy = var.task_role_policy
}

resource "aws_iam_role_policy_attachment" "task_role_attachment" {
  role       = module.ecs_alb_service_task.task_role_name
  policy_arn = aws_iam_policy.task_role_policy.arn
}

# Attach task exec roles
resource "aws_iam_policy" "task_exec_role_policy" {
  name   = "${module.label.id}-exec"
  policy = var.task_exec_role_policy
}

resource "aws_iam_role_policy_attachment" "task_exec_role_attachment" {
  role       = module.ecs_alb_service_task.task_exec_role_name
  policy_arn = aws_iam_policy.task_exec_role_policy.arn
}

# Attach security group rule to service for ingress
resource "aws_security_group_rule" "allow_tpc_ingress" {
  type              = "ingress"
  from_port         = 8080
  to_port           = 8080
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = module.ecs_alb_service_task.service_security_group_id
}

resource "aws_security_group_rule" "allow_service_ingress_to_redis" {
  count                    = var.redis_enabled ? 1 : 0
  type                     = "ingress"
  from_port                = var.redis_port
  to_port                  = var.redis_port
  protocol                 = "tcp"
  source_security_group_id = module.ecs_alb_service_task.service_security_group_id
  security_group_id        = var.redis_security_group_id
}

module "rds_instance" {
  source               = "git::https://github.com/cloudposse/terraform-aws-rds.git?ref=tags/0.26.0"
  database_name        = module.label.id
  database_user        = var.database_user
  database_password    = var.database_password
  database_port        = var.database_port
  multi_az             = var.multi_az
  storage_type         = var.storage_type
  allocated_storage    = var.allocated_storage
  storage_encrypted    = var.storage_encrypted
  engine               = var.engine
  engine_version       = var.engine_version
  instance_class       = var.instance_class
  db_parameter_group   = var.db_parameter_group
  publicly_accessible  = var.publicly_accessible
  vpc_id               = var.vpc_id
  subnet_ids           = var.private_subnet_ids
  security_group_ids   = [var.vpc_default_security_group_id]
  apply_immediately    = var.apply_immediately
  availability_zone    = var.availability_zone
}

resource "aws_security_group_rule" "allow_service_ingress_to_db" {
  type                     = "ingress"
  from_port                = var.database_port
  to_port                  = var.database_port
  protocol                 = "tcp"
  source_security_group_id = module.ecs_alb_service_task.service_security_group_id
  security_group_id        = module.rds_instance.security_group_id
}
