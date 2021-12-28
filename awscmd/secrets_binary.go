package awscmd

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type InputSecretsBinaryGet struct {
	Region string
	ID     string
}

type OutputSecretsBinaryGet struct {
	Bytes []byte
}

func SecretsBinaryGet(ctx context.Context, input *InputSecretsBinaryGet) (*OutputSecretsBinaryGet, error) {
	sess, err := NewSession(input.Region)
	if err != nil {
		return nil, err
	}

	svc := secretsmanager.New(sess, aws.NewConfig().WithRegion(input.Region))
	secretsInput := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(input.ID),
		VersionStage: aws.String("AWSCURRENT"),
	}
	result, err := svc.GetSecretValueWithContext(ctx, secretsInput)
	if err != nil {
		return nil, fmt.Errorf("Failed to get secret value: %w", err)
	}

	var out OutputSecretsBinaryGet
	if result.SecretBinary != nil {
		out.Bytes = make([]byte, len(result.SecretBinary))
		copy(out.Bytes, result.SecretBinary)
	} else {
		return nil, fmt.Errorf("Failed to fetch secret value (SecretBinary is empty): %w", err)
	}

	return &out, nil
}
