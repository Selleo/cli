package awscmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type InputSSMGetParameters struct {
	Path   string
	Region string
}

type Parameters map[string]string

type OutputSSMGetParameters struct {
	Parameters
}

func SSMGetParameters(ctx context.Context, input *InputSSMGetParameters) (*OutputSSMGetParameters, error) {
	parameters := Parameters{}
	out := &OutputSSMGetParameters{
		Parameters: parameters,
	}

	sess, err := NewSession(input.Region)
	if err != nil {
		return nil, err
	}
	svc := ssm.New(sess)
	err = svc.GetParametersByPathPagesWithContext(ctx, &ssm.GetParametersByPathInput{
		Path:           aws.String(fmt.Sprint(input.Path, "/")),
		WithDecryption: aws.Bool(true),
	}, func(o *ssm.GetParametersByPathOutput, lastPage bool) bool {
		for _, v := range o.Parameters {
			splits := strings.Split(*v.Name, "/") // split "a/b/c/ENV" by "/" to extract ENV
			env := splits[len(splits)-1]
			value := *v.Value
			parameters[env] = string(value)
		}
		return !lastPage
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}
