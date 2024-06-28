resource "aws_ecs_cluster" "chat_app_cluster" {
  name = "chat_app_cluster"
}

resource "aws_ecs_service" "app" {
  name            = "chat_app_service"
  cluster         = aws_ecs_cluster.chat_app_cluster.id
  task_definition = aws_ecs_task_definition.chat_app_task_definition.id
  launch_type     = "FARGATE"
}

resource "aws_ecs_task_definition" "chat_app_task_definition" {
  family                   = "chat_app_tasks"
  cpu                      = 256
  memory                   = 512
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  container_definitions = jsonencode([
    {
      name  = "chat_app"
      image = "tuannamnguyen290602/go-chat-app"
      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
        }
      ]
    }
  ])
}
