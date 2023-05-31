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

Create a config for AWS MFA in ~/.aws/MfaConfig where:-

* `token_arn` is the ARN of an MFA token registered in AWS to your user
* `lifetime` is the number of seconds our AWS STS token will be valid (max 3600)
* `temp_user` is a local temp user profile which will be created to hold the STS token

```
token_arn = arn:aws:iam::123456789012:mfa/my-authenticator-token
lifetime = 3600
temp_user = mfa
```

Set up AWS CLI authentication as normal in ~/.aws/credentials

```
[default]
aws_access_key_id     = AKIA1234567890ABCDEF
aws_secret_access_key = 1234567890abcdef1234567890abcdef12345678
```

Then add extra profiles in ~/.aws/credentials which depend on our temporary user profile. We do this by setting `source_profile` to the value of `temp_user` key from `~/.aws/MfaConfig`

```
[production]
role_arn = arn:aws:iam::123450000000:role/production-admin
source_profile = mfa
region = eu-west-1

[development]
role_arn = arn:aws:iam::543210000000:role/development-admin
source_profile = mfa
region = eu-west-1
```

## Usage

Execute this command to request a new STS token

```
awsmfa
```

It will modify ~/.aws/credentials to add/update the temp creds each time your run it (in our configuration, it adds the [mfa] section)

```
[default]
aws_access_key_id     = AKIA1234567890ABCDEF
aws_secret_access_key = 1234567890abcdef1234567890abcdef12345678

[mfa]
aws_access_key_id     = ASIAxxxxxxxxxxxxxxxx
aws_secret_access_key = 9U7xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
aws_session_token     = U2VjcmV0IFNlc3Npb24gVG9rZW4gLSBTZWNyZXQgU2Vzc2lvbiBUb2tlbgo=

[production]
role_arn = arn:aws:iam::123450000000:role/production-admin
source_profile = mfa
region = eu-west-1

[development]
role_arn = arn:aws:iam::543210000000:role/development-admin
source_profile = mfa
region = eu-west-1
```

Then, just execute aws cli commands as normal specifying the correct profile for the role you need

```
aws ec2 describe-instances --profile development
```

