resource "aws_ecs_cluster" "chat_app_cluster" {
  name = "chat_app_cluster"
}

resource "aws_ecs_service" "app" {
  name        = "chat_app_service"
  cluster     = aws_ecs_cluster.chat_app_cluster.id
  launch_type = "FARGATE"
}

resource "aws_ecs_task_definition" "chat_app_task" {
  family = "chat_app_tasks"
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
