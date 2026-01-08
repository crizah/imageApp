
import json
import boto3
from datetime import datetime

# AWS clients
dynamodb = boto3.client("dynamodb", region_name="us-east-1")
# kms = boto3.client("kms", region_name="us-east-1")
sns = boto3.client("sns", region_name="us-east-1")


def handler(event, context):
    print("=== POST CONFIRMATION TRIGGER ===")
    print("Full event:", json.dumps(event, indent=2))

    try:
        # Cognito trigger structure
        username = event.get("userName")
        email = event.get("request", {}).get("userAttributes", {}).get("email")

        print("Username:", username)
        print("Email:", email)

        if not username:
            raise Exception("Missing username")

        if not email:
            print("Email not found, using placeholder")

        # # 1. Create KMS key
        # print("Creating KMS key...")
        # key_response = kms.create_key(
        #     Description=f"Key for user {username}",
        #     Tags=[
        #         {
        #             "TagKey": "Username",
        #             "TagValue": username
        #         }
        #     ]
        # )

        # kms_key_id = key_response["KeyMetadata"]["KeyId"]
        # print("✅ Created KMS key:", kms_key_id)

        # Create SNS topic
        print("Creating SNS topic...")
        topic_response = sns.create_topic(
            Name=f"user-{username}-notifications"
        )

        topic_arn = topic_response["TopicArn"]
        print(" Created SNS topic:", topic_arn)

        # Subscribe email if available
        if email:
            print("Subscribing email to SNS...")
            sns.subscribe(
                TopicArn=topic_arn,
                Protocol="email",
                Endpoint=email
            )
            print(" Subscribed email to SNS")
        else:
            print(" Skipping email subscription (no email provided)")

        # 4. Save to DynamoDB
        print("Saving to DynamoDB...")
        dynamodb.put_item(
            TableName="Users",
            Item={
                "username": {"S": username},
                "email": {"S": email or "no-email-provided"},
                # "kmsKeyId": {"S": kms_key_id},
                "snsTopicArn": {"S": topic_arn},
                "createdAt": {"S": datetime.utcnow().isoformat()}
            }
        )

        print("User saved to DynamoDB")
        print("=== SUCCESS ===")

        # MUST return event to Cognito
        return event

    except Exception as e:
        print("ERROR:", str(e))
        # Do not block user signup
        return event