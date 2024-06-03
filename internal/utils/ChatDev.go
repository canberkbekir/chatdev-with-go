package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"strings"
)

type ChatDev struct {
	llm                *openai.LLM
	Config             *Config
	fullMessageHistory []string
}

func NewChatDev() (*ChatDev, error) {
	llm, err := openai.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize OpenAI client: %w", err)
	}

	config := NewConfig()

	return &ChatDev{
		llm:    llm,
		Config: config,
	}, nil
}

func (c *ChatDev) CreatePrompt(params map[string]any) (string, error) {
	promptTemplate := prompts.NewPromptTemplate(strings.Join(c.Config.ChatDev.RoleConfig.Roles["prompt engineer"], ", "), []string{"chatdev_prompt", "task"})
	formattedPrompt, err := promptTemplate.Format(params)
	if err != nil {
		return "", fmt.Errorf("failed to format prompt: %w", err)
	}

	messageHistory := []llms.MessageContent{llms.TextParts(llms.ChatMessageTypeSystem, formattedPrompt)}
	response, err := c.llm.GenerateContent(context.Background(), messageHistory)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	return strings.TrimPrefix(response.Choices[0].Content, "<INFO> "), nil
}

func (c *ChatDev) formatPrompt(agent []string, params map[string]any) (string, error) {
	promptTemplate := prompts.NewPromptTemplate(strings.Join(agent, ", "), []string{"chatdev_prompt", "task"})
	formattedPrompt, err := promptTemplate.Format(params)
	if err != nil {
		return "", fmt.Errorf("failed to format prompt: %w", err)
	}
	return formattedPrompt, nil

}

func (c *ChatDev) generateResponse(ctx context.Context, messageHistory []llms.MessageContent) (string, error) {
	response, err := c.llm.GenerateContent(ctx, messageHistory)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}
	return response.Choices[0].Content, nil
}

func (c *ChatDev) Run(userPrompt string, userAgent Role, assistantAgent Role, phase Phase, phaseParams map[string]any) (string, error) {
	InitLog()
	// Format CEO and CPO prompts
	userParams := map[string]any{"chatdev_prompt": c.Config.ChatDev.ChatChainConfig.BackgroundPrompt, "task": userPrompt}
	assistantParams := map[string]any{"chatdev_prompt": c.Config.ChatDev.ChatChainConfig.BackgroundPrompt, "task": userPrompt}

	userFullRole, err := c.formatPrompt(userAgent.RolePrompt, userParams)
	if err != nil {
		Logger.Fatal().Err(err).Msg("failed to format user prompt")
		return "", err
	}

	assistantFullRole, err := c.formatPrompt(assistantAgent.RolePrompt, assistantParams)
	if err != nil {
		Logger.Fatal().Err(err).Msg("failed to format assistant prompt")
		return "", err

	}

	phaseFullPrompt, err := c.formatPrompt(phase.Phase.PhasePrompt, phaseParams)
	if err != nil {
		Logger.Fatal().Err(err).Msg("failed to format phase prompt")
		return "", err

	}

	userMessageHistory := []llms.MessageContent{llms.TextParts(llms.ChatMessageTypeSystem, userFullRole)}
	assistantMessageHistory := []llms.MessageContent{llms.TextParts(llms.ChatMessageTypeSystem, assistantFullRole)}

	firstMessage := assistantFullRole + phaseFullPrompt
	Logger.Info().Msg(userAgent.RoleName + ": " + firstMessage)

	c.fullMessageHistory = append(c.fullMessageHistory, userAgent.RoleName+": ", firstMessage)

	assistantMessageHistory = append(assistantMessageHistory, llms.TextParts(llms.ChatMessageTypeHuman, firstMessage))

	// Generate initial response from CPO
	assistantResponseMessage, err := c.generateResponse(context.Background(), assistantMessageHistory)
	if err != nil {
		Logger.Fatal().Err(err).Msg("failed to generate initial assistant response")
		return "", err

	}
	Logger.Info().Msg(assistantAgent.RoleName + ": " + assistantResponseMessage)

	c.fullMessageHistory = append(c.fullMessageHistory, assistantAgent.RoleName+": ", assistantResponseMessage)

	maxTurn := 10
	// Exchange messages between CEO and CPO
	if phase.PhaseName == "coding" {
		maxTurn = 1
	}
	for i := 0; i < maxTurn; i++ {
		if phase.PhaseName == "coding" {
			WriteToFile(phase.PhaseName+"_History", c.fullMessageHistory)
			return strings.TrimPrefix(assistantResponseMessage, "<INFO> "), nil
		}
		userMessageHistory = append(userMessageHistory, llms.TextParts(llms.ChatMessageTypeHuman, assistantResponseMessage))

		userResponseMessage, err := c.generateResponse(context.Background(), userMessageHistory)
		if err != nil {
			Logger.Fatal().Err(err).Msg("failed to generate user response")
			return "", err
		}
		Logger.Info().Msg(userAgent.RoleName + ": " + userResponseMessage)
		c.fullMessageHistory = append(c.fullMessageHistory, userAgent.RoleName+": ", userResponseMessage)

		if strings.Contains(userResponseMessage, "<INFO>") {
			WriteToFile(phase.PhaseName+"_History", c.fullMessageHistory)
			return strings.TrimPrefix(userResponseMessage, "<INFO> "), nil
		}

		assistantMessageHistory = append(assistantMessageHistory, llms.TextParts(llms.ChatMessageTypeHuman, userResponseMessage))

		assistantResponseMessage, err = c.generateResponse(context.Background(), assistantMessageHistory)
		if err != nil {
			Logger.Fatal().Err(err).Msg("failed to generate assistant response")
			return "", err
		}
		Logger.Info().Msg(assistantAgent.RoleName + ": " + assistantResponseMessage)

		c.fullMessageHistory = append(c.fullMessageHistory, assistantAgent.RoleName+": ", assistantResponseMessage)
		if strings.Contains(assistantResponseMessage, "<INFO>") {
			WriteToFile(phase.PhaseName+"_History", c.fullMessageHistory)
			return strings.TrimPrefix(assistantResponseMessage, "<INFO> "), nil
		}

	}
	// write full message history to file
	WriteToFile(phase.PhaseName+"_History", c.fullMessageHistory)
	return "", errors.New("no response from user or assistant after 10 messages")
}
