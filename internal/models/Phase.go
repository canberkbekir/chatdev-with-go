package models

type PhaseDetails struct {
	AssistantRoleName string   `mapstructure:"assistant_role_name"`
	UserRoleName      string   `mapstructure:"user_role_name"`
	PhasePrompt       []string `mapstructure:"phase_prompt"`
}
