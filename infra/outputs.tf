output "alb_hostname" {
  value = "${aws_lb.chat_app_lb.dns_name}:8080"
}
