# Connector README

## Requirements

This connector uses your credentials stored in `.aws/credentials` file.
See [AWS documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-authentication-user.html) for more information

## Installation

To install this connector, you have two options:

1. **Copy the JSON below**.
2. **Paste the link** to the JSON configuration file into the connector's setup interface.

## JSON Configuration Example

```json
{
  "profile": "default",
  "buckets": ["bucket1", "bucket2", "bucket3"]
  "region": "us-east-1"
}
```

### Options Explained:

- profile (optional): The profile in your `.aws/credentials` file. If not given, will try to use `default`.
- region (**optional**): Region is required to access an endpoint. If not given, will try to use `default config`
- buckets (optional): List of buckets to sync.
