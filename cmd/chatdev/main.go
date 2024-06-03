package main

import (
	"chatdev-with-go/internal/utils"
	"fmt"
	"strings"
)

func main() {
	utils.InitLog()
	chatDev, err := utils.NewChatDev()
	if err != nil {
		utils.Logger.Fatal().Err(err).Msg("failed to initialize ChatDev")
		return
	}

	prompt, err := chatDev.CreatePrompt(map[string]interface{}{
		"chatdev_prompt": chatDev.Config.ApplicationConfig.ChatDev.ChatChainConfig.BackgroundPrompt,
		"task":           "Build a basic random password generator app that creates strong and customizable passwords for users.",
	})

	output, err := DemandAnalysisPhase(chatDev, prompt)
	if err != nil {
		utils.Logger.Fatal().Err(err).Msg("failed to run demand analysis phase")
		return
	}
	fmt.Print("Demand Analysis Phase Output: ")
	fmt.Println(output)

	params := map[string]interface{}{
		"task":     prompt,
		"modality": output,
		"ideas":    "",
	}
	output, err = LanguageChosePhase(chatDev, prompt, params)
	if err != nil {
		utils.Logger.Fatal().Err(err).Msg("failed to run language choice phase")
		return

	}
	fmt.Print("Language Choice Phase Output: ")
	fmt.Println(output)

	params["language"] = output
	output, err = CodingPhase(chatDev, prompt, params)
	if err != nil {
		utils.Logger.Fatal().Err(err).Msg("failed to run coding phase")
		return

	}
	fmt.Print("Coding Phase Output: ")
	fmt.Println(output)

	utils.WriteToFile("chatdev_output", []string{output})
}

func DemandAnalysisPhase(chatDev *utils.ChatDev, prompt string) (string, error) {
	roles := chatDev.Config.ChatDev.RoleConfig.Roles
	phase := utils.NewPhase("demandanalysis", chatDev.Config.ChatDev.PhaseConfig.Phases["demandanalysis"])
	ceoRole := utils.NewRole(phase.Phase.UserRoleName, roles[strings.ToLower(phase.Phase.UserRoleName)])
	cpoRole := utils.NewRole(phase.Phase.AssistantRoleName, roles[strings.ToLower(phase.Phase.AssistantRoleName)])
	run, err := chatDev.Run(prompt, *ceoRole, *cpoRole, *phase, map[string]any{"assistant_role": cpoRole.RoleName, "task": prompt})
	return run, err
}

func LanguageChosePhase(chatDev *utils.ChatDev, prompt string, params map[string]any) (string, error) {
	roles := chatDev.Config.ChatDev.RoleConfig.Roles
	phase := utils.NewPhase("languagechoice", chatDev.Config.ChatDev.PhaseConfig.Phases["languagechoose"])
	ceoRole := utils.NewRole(phase.Phase.UserRoleName, roles[strings.ToLower(phase.Phase.UserRoleName)])
	ctoRole := utils.NewRole(phase.Phase.AssistantRoleName, roles[strings.ToLower(phase.Phase.AssistantRoleName)])
	params["assistant_role"] = ctoRole.RoleName
	run, err := chatDev.Run(prompt, *ceoRole, *ctoRole, *phase, params)
	return run, err
}

func CodingPhase(chatDev *utils.ChatDev, prompt string, params map[string]any) (string, error) {
	roles := chatDev.Config.ChatDev.RoleConfig.Roles
	phase := utils.NewPhase("coding", chatDev.Config.ChatDev.PhaseConfig.Phases["coding"])
	ctoRole := utils.NewRole(phase.Phase.UserRoleName, roles[strings.ToLower(phase.Phase.UserRoleName)])
	programmerRole := utils.NewRole(phase.Phase.AssistantRoleName, roles[strings.ToLower(phase.Phase.AssistantRoleName)])
	params["assistant_role"] = programmerRole.RoleName
	run, err := chatDev.Run(prompt, *ctoRole, *programmerRole, *phase, params)
	return run, err
}
