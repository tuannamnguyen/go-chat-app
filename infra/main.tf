terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.44.0"
    }
  }

  backend "s3" {
    region = "ap-southeast-1"
    bucket = "chat-app-state-bucket"
    key    = "terraform.tfstate"
  }
}

provider "aws" {
  region = "ap-southeast-1"
}
