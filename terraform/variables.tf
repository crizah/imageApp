variable "region" {
    type = string
    description = "aws region"
  
}

variable "t1" {
    type = string
    description = "Tables name"
  
}

variable "t2" {
    type = string
    description = "Tables name"
  
}

variable "pk2"{
    type = string
    default = "messageID"
}
variable "pk1"{
    type = string
    default = "username"
}

variable "bucketName"{
    type = string
    description = "s3 bucket name"
}
