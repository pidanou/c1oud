# Connector README

This connector uses your application default credentials JSON file. If nothing is specified it will look for default location (for example `~/.config/gcloud/application_default_credentials.json`) or use the environment variable `GOOGLE_APPLICATION_CREDENTIALS`
See [Google documentation](https://cloud.google.com/docs/authentication/application-default-credentials) for more information

## Requirements

- `gcloud` CLI to authenticate to Google
- a `GCP` account and project

### Quick setup:

#### Setup GCP

On your GCP account:

- Create a project
- Activate Google Drive API

[More information](https://cloud.google.com/apis/docs/getting-started)

#### Authenticate

```bash
gcloud auth application-default login --scopes=https://www.googleapis.com/auth/drive.metadata.readonly,https://www.googleapis.com/auth/cloud-platform
gcloud auth application-default set-quota-project [project_id]
```

You should have a `application_default_credentials.json` file created.
You can know use the app using the default configuration.
You can rename and move that file to somewhere else and give the path to the connector with options (see below)

## Installation

To install this connector, you have two options:

1. **Copy the JSON below**.
2. **Paste the link** to the JSON configuration file into the connector's setup interface.

## JSON Configuration Example

```json
{
  "credentials_file": "/Users/johndoe/.config/gcloud/cred.json"
}
```

### Options Explained:

- credentials_file(optional): the path of the JSON file with your application default credentials. Will search for default locations if empty.
