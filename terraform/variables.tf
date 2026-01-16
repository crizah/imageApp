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
