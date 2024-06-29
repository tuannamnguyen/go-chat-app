resource "aws_lb" "chat_app_lb" {
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.chat_app_security_group.id]
  subnets            = [aws_subnet.chat_app_subnet.id]
}

resource "aws_lb_target_group" "chat_app_lb_target_group" {
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = aws_vpc.chat_app_vpc.id
  target_type = "ip"
}

resource "aws_lb_listener" "chat_app_lb_listener" {
  load_balancer_arn = aws_lb.chat_app_lb.arn
  port              = "80"
  protocol          = "HTTP"
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.chat_app_lb_target_group.arn
  }
}
