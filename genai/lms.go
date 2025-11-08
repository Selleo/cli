package genai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

type LMSUserInput struct {
	Audience string
	Style    string
	Topic    string
}

type LMSContentInput struct {
}

type LMSContentGenerator struct {
	*LLM
}

type LMSConfig struct {
	Dir      string       `yaml:"dir"`
	Sections []LMSSection `yaml:"sections"`
	Style    string       `yaml:"style"`
	Topic    string       `yaml:"topic"`
	Audience string       `yaml:"audience"`
}

type LMSSection struct {
	Title  string `yaml:"title"`
	Prompt string `yaml:"prompt"`
}

type LMSContent struct {
	Chapter []LMSChapter
	Image [][]byte
}

type LMSChapter struct {
	Title string
	Content string
}

func (l *LMSContentGenerator) GenerateProject(ctx context.Context, input LMSUserInput) (*LMSConfig, error) {
	id := fmt.Sprint(time.Now().Unix())
	dir := fmt.Sprintf("lms-%s", id)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create folder: %v", err)
	}

	config := &LMSConfig{
		Dir:      dir,
		Sections: []LMSSection{},
		Audience: input.Audience,
		Topic:    input.Topic,
		Style:    input.Style,
	}

	lock := sync.Mutex{}
	g := new(errgroup.Group)
	g.SetLimit(8)
	g.Go(func() error {
		out := l.LLM.GenerateText(ctx, LLMInput{
			Prompt: `Based on the STYLE and AUDIENCE, create prompt for OpenAI image generation that will give
			the best matching to the STYLE. Prompt should be generic that can later be joined with content topics.
			Output should only contain prompt, no extra information or text.
			Here is the data:
            ## Audience
			` + input.Audience + `
			## Style
			` + input.Style,
		})
		if out.Error == nil {
			lock.Lock()
			config.Style = out.Text
			lock.Unlock()
		}
		return out.Error
	})
	g.Go(func() error {
		out := l.LLM.GenerateJSON(ctx, LLMInput{
			Prompt: `Based on the TOPIC and AUDIENCE, generate a course structure for LMS platform
			that best maches target AUDIENCE and describe TOPIC in a friendly, easy to understand manner.
			Create a JSON as a list of instruction and prompts for other LLM agent.
			All the sections are just text-based lessons.
			JSON:
			{
			  "sections": [
			     {"title": "Title of the section", "prompt": "Prompt for given section"},
			     ...
			  ]
			}

			Here is the data:
            ## AUDIENCE
			` + input.Audience + `
			## TOPIC
			` + input.Topic,
		})

		decoder := json.NewDecoder(strings.NewReader(out.Text))
		var result LMSConfig
		if err := decoder.Decode(&result); err != nil {
			return err
		}
		lock.Lock()
		config.Sections = result.Sections
		if config.Sections == nil {
			return fmt.Errorf("no sections generated for course")
		}
		lock.Unlock()

		return out.Error
	})
	err = g.Wait()
	if err != nil {
		return nil, err
	}

	f, err := os.Create(path.Join(dir, "config.yml"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	defer encoder.Close()
	err = encoder.Encode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (l *LMSContentGenerator) GenerateContent(ctx context.Context, input *LMSConfig) (*LMSContent, error) {
	content := &LMSContent{}
	content.Image = make([][]byte, len(input.Sections))
	content.Chapter = make([]LMSChapter, len(input.Sections))

	lock := sync.Mutex{}
	g := new(errgroup.Group)
	g.SetLimit(8)
	for i, section := range input.Sections {
		id := fmt.Sprintf("content-%03d.md", i+1)
		g.Go(func() error {
			out := l.LLM.GenerateText(ctx, LLMInput{
				Prompt: `
				You are a LMS content creator.

				Based on the AUDIENCE and CHAPTER overview, create a single text based lesson
				about the topic that best matches audience and given chapter.
				Output should be markdown. Skip the tile in the markdown.

				Here is the data:
				## Audience
				` + input.Audience + `
				## Main topic
				` + input.Topic + `
				## Chapter title
				` + section.Title + `
				## Chapter guideline` +
				section.Prompt,
			})
			if out.Error == nil {
				lock.Lock()
				content.Chapter[i] = LMSChapter{
					Title: section.Title,
					Content: out.Text,
				}
				lock.Unlock()
				f, err := os.Create(path.Join(input.Dir, id))
				if err != nil {
					return err
				}
				defer f.Close()
				f.Write([]byte("# "+content.Chapter[i].Title))
				f.Write([]byte("\n"))
				f.Write([]byte(content.Chapter[i].Content))
			}
			return out.Error
		})
	}
	for i, section := range input.Sections {
		id := fmt.Sprintf("image-%03d.png", i+1)
		g.Go(func() error {
			out := l.LLM.GenerateImage(ctx, LLMInput{
				Prompt: `
				You are a LMS content illustrator.

				Based on the AUDIENCE, CHAPTER overview and STYLE, create an image
				about the topic that best matches audience and given chapter.

				Image should NOT contain any text.


				Here is the DATA:

				## Image STYLE:
				`+ input.Style +`
				## Audience
				` + input.Audience + `
				## Main topic
				` + input.Topic + `
				## Chapter title
				` + section.Title + `
				## Chapter guideline` +
				section.Prompt,
			})
			if out.Error == nil {
				lock.Lock()
				content.Image[i] = out.PNG
				lock.Unlock()
				f, err := os.Create(path.Join(input.Dir, id))
				if err != nil {
					return err
				}
				defer f.Close()
				f.Write(out.PNG)
			}
			return out.Error
		})
	}
	err := g.Wait()
	if err != nil {
		return nil, err
	}

	return content, nil
}
