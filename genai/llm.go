package genai

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

type LLM struct {
	token  string
	client *openai.Client
}

type LLMInput struct {
	Prompt string
}

type LLMOutput struct {
	Text  string
	PNG   []byte
	Error error
}

func NewLLM() *LLM {
	token := os.Getenv("OPENAI_API_KEY")
	if token == "" {
		panic("missing OPENAI_API_KEY")
	}
	client := openai.NewClient(token)

	return &LLM{
		token:  token,
		client: client,
	}
}

func (l *LLM) GenerateText(ctx context.Context, input LLMInput) LLMOutput {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    "user",
			Content: input.Prompt,
		},
	}
	resp, err := l.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT4oLatest,
		Messages: messages,
	})
	if err != nil {
		return LLMOutput{Error: err}
	}
	if len(resp.Choices) > 0 {
		return LLMOutput{
			Text: resp.Choices[0].Message.Content,
		}
	}

	return LLMOutput{Error: fmt.Errorf("LLM generated no output")}
}

func (l *LLM) GenerateJSON(ctx context.Context, input LLMInput) LLMOutput {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    "user",
			Content: input.Prompt,
		},
	}
	resp, err := l.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT4oLatest,
		Messages: messages,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
	})
	if err != nil {
		return LLMOutput{Error: err}
	}
	if len(resp.Choices) > 0 {
		return LLMOutput{
			Text: resp.Choices[0].Message.Content,
		}
	}

	return LLMOutput{Error: fmt.Errorf("LLM generated no output")}
}

func (l *LLM) GenerateImage(ctx context.Context, input LLMInput) LLMOutput {
	req := openai.ImageRequest{
		Prompt:            input.Prompt,
		Model:             openai.CreateImageModelGptImage1,
		N:                 1,
		Size:              openai.CreateImageSize1024x1024,
		Quality:           openai.CreateImageQualityHigh,
		OutputFormat:      openai.CreateImageOutputFormatPNG,
		Background:        openai.CreateImageBackgroundOpaque,
		OutputCompression: 100,
	}

	resp, err := l.client.CreateImage(ctx, req)
	if err != nil {
		return LLMOutput{Error: err}
	}

	if len(resp.Data) == 0 {
		return LLMOutput{Error: fmt.Errorf("no image generated")}
	}

	imgBytes, err := base64.StdEncoding.DecodeString(resp.Data[0].B64JSON)
	if err != nil {
		return LLMOutput{Error: fmt.Errorf("failed to decode image: %v", err)}
	}

	return LLMOutput{
		PNG: imgBytes,
	}
}
