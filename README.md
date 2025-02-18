# aws-sso

## Usage

> [!WARNING]
> This will be modifying your AWS credentials file. Please create a backup before using this program just in case.**

## Installation

```shell
brew install webdestroya/tap/awssso
```


## Shell Alternative
You don't even need this program, you can do this entirely within the awscli itself.

```bash

# Login to your profile and get fresh credentials
aws sso login --profile something

# Set the credentials env vars
eval $(aws configure export-credentials --profile something --format env)

# Subsequent AWS commands will use the creds above
aws sts get-caller-identity

```




