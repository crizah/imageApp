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
    from_port = 8082
    to_port = 8082
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["${var.ur_ip}/32"]  
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
              sudo apt update 
              cd /home/ubuntu
              sudo apt install git
              sudo apt install curl
              TOKEN=$(curl -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600" 2>/dev/null)
              PUBLIC_IP=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" -s http://169.254.169.254/latest/meta-data/public-ipv4)
              git clone https://github.com/crizah/imageApp-ec2.git
              cd imageApp-ec2

              cat > .env << 'ENVEOF'
              USER_POOL_CLIENT_ID=${aws_cognito_user_pool_client.client.id}
              USER_POOL_ID=${aws_cognito_user_pool.user_pool.id}
              AWS_REGION=${var.region}
              SECURE=${var.sec}
              WITH_INGRESS=${var.ingress}
              BUCKET_NAME=${var.bucketName}
              BACKEND_URL=http://$PUBLIC_IP:8082
              CLIENT_IP=http://$PUBLIC_IP:3000
              AWS_ACCESS_KEY_ID=${var.access_key_id}
              AWS_SECRET_ACCESS_KEY=${var.access_key}

              ENVEOF
              
               
               sudo apt install -y ca-certificates curl gnupg
               sudo install -m 0755 -d /etc/apt/keyrings
              curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
              sudo chmod a+r /etc/apt/keyrings/docker.gpg

              echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null



sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin


sudo usermod -aG docker ubuntu

sudo docker up --build -d

              

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
  value = "http://${aws_eip.app_eip.public_ip}:3000"
}