import { DynamoDBClient, ScanCommand } from "@aws-sdk/client-dynamodb";

const client = new DynamoDBClient({ region: "eu-north-1" });

export const handler = async (event) => {
  try {
    const params = {
      TableName: "Users",
      ProjectionExpression: "username" // only fetch the username attribute
    };

    const result = await client.send(new ScanCommand(params));

    // Extract usernames from DynamoDB items
    const usernames = result.Items?.map(item => item.username.S) || [];

    return {
      statusCode: 200,
      body: JSON.stringify({ usernames }),
      headers: { "Access-Control-Allow-Origin": "*" }


      
    };

  } catch (err) {
    console.error(err);
    return {
      statusCode: 500,
      body: JSON.stringify({ message: "Internal Server Error", error: err.message }),
      headers: { "Access-Control-Allow-Origin": "*" }
    };
  }
};
