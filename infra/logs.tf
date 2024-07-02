# Set up CloudWatch group and log stream and retain logs for 30 days
resource "aws_cloudwatch_log_group" "chat_app_log_group" {
  name              = "/ecs/chat_app"
  retention_in_days = 30
}
