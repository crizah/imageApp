import json
import boto3

dynamodb = boto3.client("dynamodb", region_name="us-east-1")

def handler(event, context):
    try:
        params = {
            "TableName": "Users",
            "ProjectionExpression": "username"
        }

        result = dynamodb.scan(**params)

        # Extract usernames
        usernames = [
            item["username"]["S"]
            for item in result.get("Items", [])
            if "username" in item
        ]

        return {
            "statusCode": 200,
            "headers": {
                "Access-Control-Allow-Origin": "*",
                "Access-Control-Allow-Headers": "Content-Type",
                "Access-Control-Allow-Methods": "GET,OPTIONS"
            },
            "body": json.dumps({"usernames": usernames})
        }

    except Exception as err:
        print("ERROR:", err)

        return {
            "statusCode": 500,
            "headers": {
                "Access-Control-Allow-Origin": "*",
                "Access-Control-Allow-Headers": "Content-Type",
                "Access-Control-Allow-Methods": "GET,OPTIONS"
            },
            "body": json.dumps({
                "message": "Internal Server Error",
                "error": str(err)
            })
        }
