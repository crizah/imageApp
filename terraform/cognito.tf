resource "aws_cognito_user_pool" "user_pool" {
  name = "user-pool"

#   username_attributes = ["email"] // this does that anyways 
#   alias_attributes = ["preferred_username"] // allows foir login with username as well as email
   alias_attributes = ["email", "preferred_username"]
  auto_verified_attributes = ["email"]

  password_policy {
    minimum_length = 6
  }

  verification_message_template {
    default_email_option = "CONFIRM_WITH_CODE" # confirm with link needs a user domain to be associated. stupid 
    email_subject = "Account Confirmation"
    email_message = "This is your verification code: {####}"
  }
#   email_configuration {
#     email_sending_account = "COGNITO_DEFAULT" 
#   }

  schema {
    attribute_data_type      = "String"
    developer_only_attribute = false
    mutable                  = true
    name                     = "email"
    required                 = true

    string_attribute_constraints {
      min_length = 1
      max_length = 256
    }


    
  }

  schema {
    attribute_data_type      = "String"
    developer_only_attribute = false
    mutable                  = true  
    name                     = "name"
    required                 = true

  }


  # trigger a lamda function upon sign up and after user is verified  
    lambda_config {
    post_confirmation = module.lambda_function_existing_package_local.lambda_function_arn
    }

    # Account recovery
  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_email"
      priority = 1
    }
  }
    
}

resource "aws_cognito_user_pool_client" "client" {
  name = "cognito-client"

  user_pool_id = aws_cognito_user_pool.user_pool.id
  generate_secret = false
  refresh_token_validity = 90
  prevent_user_existence_errors = "ENABLED"
  explicit_auth_flows = [
    "ALLOW_REFRESH_TOKEN_AUTH",
    "ALLOW_USER_PASSWORD_AUTH",
    "ALLOW_ADMIN_USER_PASSWORD_AUTH"
  ]


  # Token validity
  access_token_validity  = 60   # minutes
  id_token_validity      = 60   # minutes
  token_validity_units {
    access_token  = "minutes"
    id_token      = "minutes"
    refresh_token = "days"
  }


  # Read/write attributes
  read_attributes = [
    "email",
    "email_verified",
    "name"
  ]
  
  write_attributes = [
    "email",
    "name"
  ]

}


# resource "aws_cognito_user_pool_domain" "cognito-domain" {
#   domain       = "meow.auth.${var.region}.amazoncognito.com"
#   user_pool_id = "${aws_cognito_user_pool.user_pool.id}"
# }


output "user_pool_id" {
  value = aws_cognito_user_pool.user_pool.id
}

output "user_pool_client_id" {
  value = aws_cognito_user_pool_client.client.id
}

output "user_pool_endpoint" {
  value = aws_cognito_user_pool.user_pool.endpoint
}





