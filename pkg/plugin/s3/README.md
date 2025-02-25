# Plugin README

## Installation

To install this plugin, you have two options:

1. **Copy the JSON below**.
2. **Paste the link** to the JSON configuration file into the plugin's setup interface.

## JSON Configuration Example

```json
{
  "profile": "default",
  "buckets": ["bucket1", "bucket2", "bucket3"]
}
```

### Options Explained:

- profile: the profile in your `.aws/credentials` file. If not given, will try to use `default`.
- buckets: List of buckets to sync. Default to all buckets of the account.
