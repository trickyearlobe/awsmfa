package main

import (
	"fmt"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"gopkg.in/ini.v1"
)

func main() {
	// Get the awsmfa configuration from ~/.aws/MfaConfig
	homedir, _ := os.UserHomeDir()
	configFile := path.Join(homedir, ".aws", "MfaConfig")
	mfaConfig, err := ini.Load(configFile)

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(124)
	}

	// Validate the config
	token_arn := mfaConfig.Section("").Key("token_arn").String()
	temp_user := mfaConfig.Section("").Key("temp_user").String()
	lifetime := mfaConfig.Section("").Key("lifetime").MustInt64()
	if (token_arn == "") || (temp_user == "") || (lifetime < 900) || (lifetime > 3600) {
		fmt.Printf(
			"\nPlease edit %v to include the following lines\n\n"+
				"  token_arn = arn:aws:iam::<aws account id>:mfa/<token serial>\n"+
				"  temp_user = <temporary username>\n"+
				"  lifetime = <lifetime in seconds between 900 and 3600>\n\n", configFile)
		os.Exit(125)
	}

	// Start an AWS session using the default profile/creds
	stsSvc := sts.New(session.New())

	// Build a GetSessionToken Input with
	params := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(lifetime),
		SerialNumber:    aws.String(token_arn),
		TokenCode:       getPasscode(),
	}

	// Fetch the Session Token if we can
	creds, err := stsSvc.GetSessionToken(params)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(123)
	}

	// Read the credentials file
	credentialsFile := path.Join(homedir, ".aws", "credentials")
	credentials, err := ini.Load(credentialsFile)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(122)
	}

	// Write the new session token to the credentials file temporary user
	credentials.Section(temp_user).Key("aws_access_key_id").SetValue(*creds.Credentials.AccessKeyId)
	credentials.Section(temp_user).Key("aws_secret_access_key").SetValue(*creds.Credentials.SecretAccessKey)
	credentials.Section(temp_user).Key("aws_session_token").SetValue(*creds.Credentials.SessionToken)
	err = credentials.SaveTo(credentialsFile)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(121)
	}
	fmt.Printf("Your temporary token will expire at %v\n", creds.Credentials.Expiration)
}

func getPasscode() *string {
	var passcode string
	fmt.Print("Enter passcode: ")
	fmt.Scanln(&passcode)
	return aws.String(passcode)
}
