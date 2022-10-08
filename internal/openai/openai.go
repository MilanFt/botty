package openai

import (
	"context"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
)

// New returns an instance of gpt3.Client with the default engine (text-davinci-002).
func New(key string) gpt3.Client {
	return gpt3.NewClient(
		key,
		gpt3.WithDefaultEngine("text-davinci-002"),
	)
}

// NewWithEngine returns a gpt3.Client instance with the specified engine type.
func NewWithEngine(key string, engine string) gpt3.Client {
	return gpt3.NewClient(
		key,
		gpt3.WithDefaultEngine(engine),
	)
}

// Complete returns the result of the completion API call from the provided prompt and users.
func Complete(c gpt3.Client, prompt string, users []string) (string, error) {
	var stops []string
	for _, v := range users {
		stops = append(stops, v+":")
	}
	resp, err := c.Completion(context.Background(), gpt3.CompletionRequest{
		Prompt:           []string{prompt},
		MaxTokens:        gpt3.IntPtr(150),
		Temperature:      gpt3.Float32Ptr(0.9),
		TopP:             gpt3.Float32Ptr(1),
		FrequencyPenalty: 0,
		PresencePenalty:  0.6,
		Stop:             stops,
	})
	if err != nil {
		return "", err
	}
	res := strings.TrimSpace(resp.Choices[0].Text)
	return res, nil
}

// CreatePrompt returns a prompt string usable in the 'Complete' function
// from the bot's name, identity and the provided chat logs.
func CreatePrompt(name string, identity string, logs []string) string {
	var prompt string
	prompt += identity + "\n\n"
	for _, v := range logs {
		prompt += v + "\n"
	}
	prompt += name + ":"
	return prompt
}
