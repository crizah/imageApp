variable "region" {
    type = string
    description = "aws region"
  
}

variable "users" {
    type = string
    description = "Tables name"
  
}

variable "msgs" {
    type = string
    description = "Tables name"
  
}

variable "msgspk"{
    type = string
    default = "messageID"
}
variable "userspk"{
    type = string
    default = "username"
}

variable "msgsgsi"{
    type = string
    default = "recipientIndex"
}

variable "bucketName"{
    type = string
    description = "s3 bucket name"
}

variable "sec"{
    type = string
    description = "secure for https"
    
}

variable "ingress"{
    type = string
    description = "with ingress for k8s"
}

variable "ssh_key"{
    type = string
    description = "ssh key for ec2"
}

variable "ur_ip"{
    type = string
    description = "ur ip address"
}

variable "access_key_id" {
    type= string
    description = "aws access key id"
  
}

variable "access_key" {
    type = string
    description = "aws access key"
  
}
