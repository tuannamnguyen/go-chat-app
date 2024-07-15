terraform {


  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.44.0"
    }

    porkbun = {
      source  = "kyswtn/porkbun"
      version = "0.1.2"
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

provider "aws" {
  alias  = "us-east-1"
  region = "us-east-1"
}

variable "porkbun_api_key" {
  type = string
}

variable "porkbun_secret_api_key" {
  type = string
}

provider "porkbun" {
  api_key        = var.porkbun_api_key
  secret_api_key = var.porkbun_secret_api_key
}
