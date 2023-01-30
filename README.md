# AWS MFA

A basic command line client that assist with assuming roles that require MFA when using standard AWS API tokens that don't include MFA.

It aquires STS temporary credentials which it writes to an alternate user in the ~/.aws/credentials file.

The alternate user can be used by profiles in the ~/.aws/config file to force the AWS CLI to assume a new role.

## Installation

```
git clone https://github.com/trickyearlobe/awsmfa.git
cd awsmfa

# Either build the binary in the local directory
go build

# Or install the binary into your GOLANG bin directory
go install
```

## Configuration

Set up AWS CLI authentication as normal in ~/.aws/credentials

```
[default]
aws_access_key_id     = AKIA1234567890ABCDEF
aws_secret_access_key = 1234567890abcdef1234567890abcdef12345678
```

Create a config for AWS MFA in ~/.aws/MfaConfig

```
token_arn = arn:aws:iam::123456789012:mfa/my-authenticator-token
lifetime = 3600
temp_user = tokenuser
```

Finally, set up profiles in ~/.aws/config using `source_profile` to use the temporary alternate credentials

```
[default]
region = eu-west-1

[profile production]
role_arn = arn:aws:iam::123450000000:role/production-admin
region = eu-west-1
source_profile = tokenuser

[profile development]
role_arn = arn:aws:iam::543210000000:role/development-admin
source_profile = tokenuser
```

## Usage

```
awsmfa
```
