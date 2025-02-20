# aws-sso

Writes credentials to your `~/.aws/credentials` file that were created via SSO configurations. This allows you to run old applications that rely on AccessKey/SecretKey credentials while still allowing you to get those credentials via SSO.


## Installation

```shell
brew install webdestroya/tap/awssso
```

## Usage

> [!WARNING]
> This will be modifying your AWS credentials file. Please create a backup before using this program just in case.

Pull credentials for the listed profiles and update the credentials file
```shell
awssso sync profile1 profile2 profile3
```


## Shell Alternative
You don't even need this program, you can do this entirely within the awscli itself. (Assuming your program can use environment variables instead of aws credentials file.)

```bash

# Login to your profile and get fresh credentials
aws sso login --profile something

# Set the credentials env vars
eval $(aws configure export-credentials --profile something --format env)

# Subsequent AWS commands will use the creds above
aws sts get-caller-identity

```




