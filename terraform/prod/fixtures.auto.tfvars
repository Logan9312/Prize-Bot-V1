region = "us-west-2"

namespace = "bot"

stage = "prod"

name = "auction"

attributes = []

delimiter = "-"

tags = {}

vpc_id = "vpc-0c13caa4c33e809b8"

vpc_default_security_group_id = "sg-09eaea89b68deebb7"

ecs_cluster_arn = "arn:aws:ecs:us-west-2:357595321916:cluster/disc-prod-bots"

private_subnet_ids = [
  "subnet-0bd270e58f03d7844",
  "subnet-008e39bacbb19d2e6",
]

load_balancer_listener_arn = "arn:aws:elasticloadbalancing:us-west-2:357595321916:listener/app/disc-prod-bots/780f559464828620/18b256984fa9b334"

load_balancer_listener_paths = ["/auction-bot/*"]

target_group_port = 80

target_group_protocol = "HTTP"

target_group_target_type = "ip"

health_check_path = "/auction-bot/status"

health_check_timeout = 10

health_check_healthy_threshold = 2

health_check_unhealthy_threshold = 2

health_check_interval = 15

health_check_matcher = "200"

container_image = "357595321916.dkr.ecr.us-west-2.amazonaws.com/auction-bot:prod"

container_memory = 512

container_memory_reservation = 450

container_port_mappings = [
  {
    containerPort = 8080
    hostPort      = 8080
    protocol      = "tcp"
  }
]

container_port = 8080

container_cpu = 256

container_essential = true

container_environment = [
  {
    name  = "LISTENER_PORT",
    value = "8080"
  }
]

container_secrets = [
  {
    name      = "ENV_VARS"
    valueFrom = "arn:aws:secretsmanager:us-west-2:357595321916:secret:prod/auction-bot-Nmd95V"
  }
]

container_readonly_root_filesystem = false

ecs_launch_type = "FARGATE"

ignore_changes_task_definition = true

network_mode = "null"

assign_public_ip = false

propagate_tags = "TASK_DEFINITION"

deployment_minimum_healthy_percent = 0

deployment_maximum_percent = 100

deployment_controller_type = "ECS"

desired_count = 1

task_memory = 512

task_cpu = 256

task_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["secretsmanager:GetSecretValue"],
      "Resource": "*"
    }
  ]
}
EOF

task_exec_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["secretsmanager:GetSecretValue"],
      "Resource": "*"
    }
  ]
}
EOF


deletion_protection = false

database_port = 3306

multi_az = false

storage_type = "gp2"

storage_encrypted = false

allocated_storage = 5

engine = "postgres"

engine_version = "12.5"

major_engine_version = "12"

instance_class = "db.t3.micro"

db_parameter_group = "postgres12"

publicly_accessible = true

apply_immediately = true
