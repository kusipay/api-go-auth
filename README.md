# Getting started

Template code for initializing serverless with go1.20.

## Initialization

- Copy paste the files into a new folder.
- Change the name of the service (for serverless repo use the format "sls-{name}").
- Change the name in go mod.
- Push the code to the remote repo.

## Developing

Run the following commands to download the necesary packages

```bash
go mod download
go mod tidy
go mod verify
```
Remove/Add the necessary packages for the service.

## Content

In this template theres 5 lambdas, 2 http api, 2 sqs and 1 dynamo table.
