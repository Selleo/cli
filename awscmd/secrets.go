package awscmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type InputSecretsAll struct {
	Region string
	ID     string
}

type OutputSecretsAll struct {
	Secrets map[string]string
}

func SecretsAll(ctx context.Context, input *InputSecretsAll) (*OutputSecretsAll, error) {
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

	var out OutputSecretsAll
	if result.SecretString != nil {
		err = json.NewDecoder(strings.NewReader(*result.SecretString)).Decode(&out.Secrets)
		if err != nil {
			return nil, fmt.Errorf("Failed to decode secrets into key-value map: %w", err)
		}
	} else {
		return nil, fmt.Errorf("Failed to fetch secret value (SecretString is empty): %w", err)
	}

	return &out, nil
}

type InputSecretsPut struct {
	Region     string
	ID         string
	NewSecrets map[string]string
}

type OutputSecretsPut struct{}

func SecretsPut(ctx context.Context, input *InputSecretsPut) (*OutputSecretsPut, error) {
	all, err := SecretsAll(ctx, &InputSecretsAll{
		Region: input.Region,
		ID:     input.ID,
	})
	if err != nil {
		return nil, err
	}

	for k, v := range input.NewSecrets {
		all.Secrets[k] = v
	}

	secretBytes, err := json.Marshal(all.Secrets)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshall secrets to json: %w", err)
	}

	sess, err := NewSession(input.Region)
	if err != nil {
		return nil, err
	}

	secretsString := string(secretBytes)
	svc := secretsmanager.New(sess, aws.NewConfig().WithRegion(input.Region))
	secretsInput := &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(input.ID),
		SecretString: &secretsString,
	}
	_, err = svc.PutSecretValueWithContext(ctx, secretsInput)
	if err != nil {
		return nil, fmt.Errorf("Failed to set secrets: %w", err)
	}

	return &OutputSecretsPut{}, nil
}
