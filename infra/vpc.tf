resource "aws_vpc" "chat_app_vpc" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_internet_gateway" "chat_app_internet_gateway" {
  vpc_id = aws_vpc.chat_app_vpc.id
}

resource "aws_route_table" "chat_app_vpc_route_table" {
  vpc_id = aws_vpc.chat_app_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.chat_app_internet_gateway.id
  }
}

# Route subnet to the internet
resource "aws_route_table_association" "chat_app_table_association" {
  route_table_id = aws_route_table.chat_app_vpc_route_table.id
  subnet_id      = aws_subnet.chat_app_subnet.id
}

# Create public subnets
resource "aws_subnet" "chat_app_subnet" {
  vpc_id                  = aws_vpc.chat_app_vpc.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  availability_zone       = "ap-southeast-1a"
}

# Allow traffic to port 8080
resource "aws_security_group" "chat_app_security_group" {
  vpc_id = aws_vpc.chat_app_vpc.id
}

resource "aws_vpc_security_group_ingress_rule" "chat_app_ingress_rule" {
  security_group_id = aws_security_group.chat_app_security_group.id

  ip_protocol = "tcp"
  cidr_ipv4   = "0.0.0.0/0"
  from_port   = 8080
  to_port     = 8080
}

resource "aws_vpc_security_group_egress_rule" "chat_app_egress_rule" {
  security_group_id = aws_security_group.chat_app_security_group.id

  ip_protocol = "-1"
  cidr_ipv4   = "0.0.0.0/0"
}
