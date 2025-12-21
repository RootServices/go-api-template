# Environment File

environment variables are defined in a .env file. They are located in the [/env](../env) directory.

## Usage

### Local Development


| Variable | Description | 
|----------|-------------|
|ENV|Must be set to `local`|
|DB_HOST|The host of the database.|
|DB_PORT|The port of the database.|
|DB_NAME|The name of the database.|
|DB_USER|The user of the database.|
|DB_PASSWORD|The password of the database.|
|DB_SSL_MODE|The SSL mode of the database.|


### GCP Cloud Run

| Variable | Description | 
|----------|-------------|
| Variable | Description | 
|----------|-------------|
|ENV|The environment the application is running in.|
|GCP_PROJECT_NUMBER|The GCP Project Number|
|GCP_PROJECT_ID|The GCP Project ID|
|DB_HOST|The host of the database.|
|DB_PORT|The port of the database.|
|DB_NAME|The name of the database.|
|DB_USER_KEY|The key in secret manager for the database user.|
|DB_PASSWORD_KEY|The key in secret manager for the database password.|
|DB_SSL_MODE|The SSL mode of the database.|