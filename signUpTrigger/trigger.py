import json
import boto3
import os
from datetime import datetime

# from env variables set by terraform
REGION = os.environ.get("AWS_REGION_NAME", "us-east-1") # defaults JUST IN CASE
DYNAMODB_TABLE = os.environ.get("DYNAMODB_TABLE", "Users")

# server
dynamodb = boto3.client("dynamodb", region_name=REGION)
sns = boto3.client("sns", region_name=REGION)

def handler(event, context):
    
    try:
        # get user info from event
        username = event.get("userName")
        user_attributes = event.get("request", {}).get("userAttributes", {})
        email = user_attributes.get("email")
        user_sub = user_attributes.get("sub") 
        
      
        # need all 3
        if not username:
            raise ValueError("Missing username in Cognito event")
        if not email:
            raise ValueError("Missing email in Cognito event")
        if not user_sub:
            raise ValueError("Missing user sub (ID) in Cognito event")
        
        # create SNS topic for user notifications
       
        topic_name = f"user-{username}-notifications"
        topic_response = sns.create_topic(Name=topic_name)
        topic_arn = topic_response["TopicArn"]
     
        
        # Subscribe user's email to their topic
       
        sns.subscribe(
            TopicArn=topic_arn,
            Protocol="email",
            Endpoint=email
        )
      
        
        # Save user to DynamoDB
        
        dynamodb.put_item(
            TableName=DYNAMODB_TABLE,
            Item={
                "username": {"S": username},
                # "userId": {"S": user_sub},  # Use Cognito sub as primary ID
                "email": {"S": email},
                "snsTopicArn": {"S": topic_arn},
                "emailVerified": {"BOOL": True}, 
                "createdAt": {"S": datetime.now().isoformat()},
                # "updatedAt": {"S": datetime.utcnow().isoformat()}
            }
        )
        print(" User saved to DynamoDB")
        

        
        # return the event to Cognito
        return event
        
    except Exception as e:
        print(f" ERROR: {str(e)}")
        print(f"Error type: {type(e).__name__}")
        
        # Still return event so Cognito doesn't fail the confirmation
        # The user is already verified, we just couldn't do post-processing
        return event