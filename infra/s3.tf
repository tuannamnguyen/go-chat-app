resource "aws_s3_bucket" "chat_app_env_bucket" {
  force_destroy = true
  bucket        = "chat-app-env-bucket"
}

resource "aws_s3_object" "env_file" {
  bucket = aws_s3_bucket.chat_app_env_bucket.id
  key    = ".env"
  source = "../.env"
}
