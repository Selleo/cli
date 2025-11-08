package generators

import (
	"context"
	"os"

	"github.com/Selleo/cli/genai"
)

type LLM struct {}

func (l *LLM) Text(ctx context.Context, input genai.LLMInput) genai.LLMOutput {
	return genai.NewLLM().GenerateText(ctx, input)
}

func (l *LLM) Image(ctx context.Context, out string, input genai.LLMInput) error {
	res := genai.NewLLM().GenerateImage(ctx, input)
	if res.Error != nil {
		return res.Error
	}

	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()


	_, err = f.Write(res.PNG)
	if err != nil {
		return err
	}

	return nil
}
