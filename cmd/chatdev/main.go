package main

import (
	"chatdev-with-go/internal/config"
	"chatdev-with-go/internal/utils"
	"context"
	"fmt"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"strings"
)

func main() {
	testPrompt := "Write love 10 times with golang"
	utils.InitLog()
	utils.Logger.Info().Msg("Starting ChatDev with Go")
	test := config.NewConfig()
	utils.Logger.Info().Msg("ChatDev with Go started successfully")
	llm, err := openai.New()
	if err != nil {
		utils.Logger.Fatal().Err(err)
	}
	result := strings.Join(test.ApplicationConfig.ChatDev.RoleConfig.Programmer, " ")
	prompt := prompts.NewPromptTemplate(result, []string{"task", "chatdev_prompt"})
	llmChain := chains.NewLLMChain(llm, prompt)
	ctx := context.Background()
	out, err := chains.Call(ctx, llmChain, map[string]any{
		"task":           testPrompt,
		"chatdev_prompt": testPrompt,
	})
	if err != nil {
		utils.Logger.Fatal().Err(err)
	}
	fmt.Println(out)

}
