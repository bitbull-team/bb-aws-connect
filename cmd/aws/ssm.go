package aws

import (
	"awslib"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/urfave/cli/v2"
)

// SSMListAvailableInstances return a list of available instaces
func SSMListAvailableInstances(c *cli.Context) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           c.String("profile"),
		SharedConfigState: session.SharedConfigEnable,
	}))

	_, err := awslib.EC2ListInstances(sess, []awslib.TagFilter{
		awslib.TagFilter{
			Name:  "Environment",
			Value: c.String("env"),
		},
	})

	if err != nil {
		return cli.Exit("Error during EC2 instance list: "+err.Error(), -1)
	}

	return nil
}
