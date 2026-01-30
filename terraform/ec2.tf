

# IAM Role for EC2 to access AWS services
resource "aws_iam_role" "ec2_role" {
  name = "chat-app-ec2-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
  })
}

# allowing cognito and dynamo access
resource "aws_iam_role_policy" "ec2_policy" {
  name = "chat-app-ec2-policy"
  role = aws_iam_role.ec2_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "cognito-idp:*",
          "dynamodb:*"
        ]
        Resource = "*"
      }
    ]
  })
}

# attatch to ec2
resource "aws_iam_instance_profile" "ec2_profile" {
  name = "chat-app-ec2-profile"
  role = aws_iam_role.ec2_role.name
}

# security groups
resource "aws_security_group" "app_sg" {
  name        = "chat-app-sg"
  description = "Allow HTTP, HTTPS, SSH"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

   ingress {
    from_port   = 3000
    to_port     = 3000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["${var.ur_ip}/32"]  # SSH to your ip
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}


resource "aws_instance" "app_server" {
  ami                    = "ami-0b6c6ebed2801a5cb"
  instance_type          = "t3.micro" 
  iam_instance_profile   = aws_iam_instance_profile.ec2_profile.name
  vpc_security_group_ids = [aws_security_group.app_sg.id]
  key_name               = var.ssh_key  

  user_data = <<-EOF
                               
              #!/bin/bash
              set -e
              # install docker and coker compose
              sudo apt-get update
              sudo apt-get upgrade

              sudo apt-get install \
                    ca-certificates \
                    curl \
                    gnupg \
                    lsb-release \
                    git
              sudo apt install apt-transport-https ca-certificates curl software-properties-common

              curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

              sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable"

              sudo apt update

              sudo apt install docker-ce

        

              # docker compose

              sudo curl -L "https://github.com/docker/compose/releases/download/1.28.5/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
              sudo chmod +x /usr/local/bin/docker-compose


              cd /home/ubuntu

              # clone

              git clone  https://github.com/crizah/imageApp.git
              cd imageApp

              # create the .env file

               
              
              cat > .env << 'ENVEOF'
              USER_POOL_CLIENT_ID=${aws_cognito_user_pool_client.client.id}
              USER_POOL_ID=${aws_cognito_user_pool.user_pool.id}
              AWS_REGION=${var.region}
              SECURE=${var.sec}
              WITH_INGRESS=${var.ingress}
              BUCKET_NAME=${var.bucketName}
              BACKEND_URL=http://backend:8082
              ENVEOF

              # run docker compose in backgroung

              docker compose up --build
              EOF



  tags = {
    Name = "imageApp"
  }
}


resource "aws_eip" "app_eip" {
  instance = aws_instance.app_server.id
}

output "public_ip" {
  value = aws_eip.app_eip.public_ip
}

output "app_url" {
  value = "http://${aws_eip.app_eip.public_ip}"
}