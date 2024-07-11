variable "image_tag" {
  type    = string
  default = "latest"
}

resource "aws_ecs_cluster" "chat_app_cluster" {
  name = "chat_app_cluster"
}

resource "aws_ecs_service" "chat_app_service" {
  name                 = "chat_app_service"
  cluster              = aws_ecs_cluster.chat_app_cluster.id
  task_definition      = aws_ecs_task_definition.chat_app_task_definition.id
  launch_type          = "FARGATE"
  force_new_deployment = true
  desired_count        = 1
  network_configuration {
    subnets          = [aws_subnet.chat_app_subnet.id, aws_subnet.chat_app_subnet_2.id]
    security_groups  = [aws_security_group.chat_app_security_group.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.chat_app_lb_target_group.id
    container_name   = "chat_app"
    container_port   = 8080
  }

  depends_on = [aws_lb_listener.chat_app_lb_listener]
}

resource "aws_ecs_task_definition" "chat_app_task_definition" {
  family                   = "chat_app_tasks"
  cpu                      = 256
  memory                   = 512
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  container_definitions = jsonencode([
    {
      name  = "chat_app"
      image = "tuannamnguyen290602/go-chat-app:${var.image_tag}"
      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
        }
      ]
      dependsOn = [
        {
          containerName = "redis"
          condition     = "START"
        }
      ]
      environmentFiles = [
        {
          type  = "s3"
          value = aws_s3_object.env_file.arn
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.chat_app_log_group.name
          awslogs-region        = "ap-southeast-1"
          awslogs-stream-prefix = "ecs_chatapp"
        }
      }
    },
    {
      name      = "redis"
      image     = "redis"
      essential = true
      portMappings = [
        {
          containerPort = 6379
          hostPort      = 6379
          protocol      = "tcp"
        }
      ]
      command = [
        "redis-server",
        "--save",
        "60",
        "1",
        "--loglevel",
        "warning"
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.chat_app_log_group.name
          awslogs-region        = "ap-southeast-1"
          awslogs-stream-prefix = "ecs_redis"
        }
      }
    }
  ])
  depends_on = [aws_s3_object.env_file]
}
