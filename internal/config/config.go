package config

import (
	"chatdev-with-go/internal/utils"
	"github.com/spf13/viper"
)

type Config struct {
	*ApplicationConfig
}

type ApplicationConfig struct {
	ChatDev ChatDevConfig `mapstructure:"chat_dev"`
	Server  ServerConfig  `mapstructure:"server"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type ChatDevConfig struct {
	RoleConfig      RoleConfig      `mapstructure:"role_config"`
	PhaseConfig     PhaseConfig     `mapstructure:"phase_config"`
	ChatChainConfig ChatChainConfig `mapstructure:"chat_chain_config"`
}
type ChatChainConfig struct {
	Chain []struct {
		Phase       string `mapstructure:"phase"`
		PhaseType   string `mapstructure:"phaseType"`
		MaxTurnStep int    `mapstructure:"max_turn_step,omitempty"`
		NeedReflect bool   `mapstructure:"need_reflect,omitempty"`
		CycleNum    int    `mapstructure:"cycleNum,omitempty"`
		Composition []struct {
			Phase       string `mapstructure:"phase"`
			PhaseType   string `mapstructure:"phaseType"`
			MaxTurnStep int    `mapstructure:"max_turn_step"`
			NeedReflect bool   `mapstructure:"need_reflect"`
		} `mapstructure:"Composition,omitempty"`
	} `mapstructure:"chain"`
	Recruitments       []string `mapstructure:"recruitments"`
	ClearStructure     bool     `mapstructure:"clear_structure"`
	GuiDesign          bool     `mapstructure:"gui_design"`
	GitManagement      bool     `mapstructure:"git_management"`
	WebSpider          bool     `mapstructure:"web_spider"`
	SelfImprove        bool     `mapstructure:"self_improve"`
	IncrementalDevelop bool     `mapstructure:"incremental_develop"`
	WithMemory         bool     `mapstructure:"with_memory"`
	BackgroundPrompt   string   `mapstructure:"background_prompt"`
}

type RoleConfig struct {
	ChiefExecutiveOfficer     []string `mapstructure:"Chief Executive Officer"`
	ChiefProductOfficer       []string `mapstructure:"Chief Product Officer"`
	Counselor                 []string `mapstructure:"Counselor"`
	ChiefTechnologyOfficer    []string `mapstructure:"Chief Technology Officer"`
	ChiefHumanResourceOfficer []string `mapstructure:"Chief Human Resource Officer"`
	Programmer                []string `mapstructure:"Programmer"`
	CodeReviewer              []string `mapstructure:"Code Reviewer"`
	SoftwareTestEngineer      []string `mapstructure:"Software Test Engineer"`
	ChiefCreativeOfficer      []string `mapstructure:"Chief Creative Officer"`
}

type PhaseConfig struct {
	DemandAnalysis struct {
		AssistantRoleName string   `mapstructure:"assistant_role_name"`
		UserRoleName      string   `mapstructure:"user_role_name"`
		PhasePrompt       []string `mapstructure:"phase_prompt"`
	} `mapstructure:"DemandAnalysis"`
	LanguageChoose struct {
		AssistantRoleName string   `mapstructure:"assistant_role_name"`
		UserRoleName      string   `mapstructure:"user_role_name"`
		PhasePrompt       []string `mapstructure:"phase_prompt"`
	} `mapstructure:"LanguageChoose"`
	Coding struct {
		AssistantRoleName string   `mapstructure:"assistant_role_name"`
		UserRoleName      string   `mapstructure:"user_role_name"`
		PhasePrompt       []string `mapstructure:"phase_prompt"`
	} `mapstructure:"Coding"`
	ArtDesign struct {
		AssistantRoleName string   `mapstructure:"assistant_role_name"`
		UserRoleName      string   `mapstructure:"user_role_name"`
		PhasePrompt       []string `mapstructure:"phase_prompt"`
	} `mapstructure:"ArtDesign"`
	ArtIntegration struct {
		AssistantRoleName string   `mapstructure:"assistant_role_name"`
		UserRoleName      string   `mapstructure:"user_role_name"`
		PhasePrompt       []string `mapstructure:"phase_prompt"`
	} `mapstructure:"ArtIntegration"`
	CodeComplete struct {
		AssistantRoleName string   `mapstructure:"assistant_role_name"`
		UserRoleName      string   `mapstructure:"user_role_name"`
		PhasePrompt       []string `mapstructure:"phase_prompt"`
	} `mapstructure:"CodeComplete"`
	CodeReviewComment struct {
		AssistantRoleName string   `mapstructure:"assistant_role_name"`
		UserRoleName      string   `mapstructure:"user_role_name"`
		PhasePrompt       []string `mapstructure:"phase_prompt"`
	} `mapstructure:"CodeReviewComment"`
	CodeReviewModification struct {
		AssistantRoleName string   `mapstructure:"assistant_role_name"`
		UserRoleName      string   `mapstructure:"user_role_name"`
		PhasePrompt       []string `mapstructure:"phase_prompt"`
	} `mapstructure:"CodeReviewModification"`
	TestErrorSummary struct {
		AssistantRoleName string   `mapstructure:"assistant_role_name"`
		UserRoleName      string   `mapstructure:"user_role_name"`
		PhasePrompt       []string `mapstructure:"phase_prompt"`
	} `mapstructure:"TestErrorSummary"`
	TestModification struct {
		AssistantRoleName string   `mapstructure:"assistant_role_name"`
		UserRoleName      string   `mapstructure:"user_role_name"`
		PhasePrompt       []string `mapstructure:"phase_prompt"`
	} `mapstructure:"TestModification"`
	EnvironmentDoc struct {
		AssistantRoleName string   `mapstructure:"assistant_role_name"`
		UserRoleName      string   `mapstructure:"user_role_name"`
		PhasePrompt       []string `mapstructure:"phase_prompt"`
	} `mapstructure:"EnvironmentDoc"`
	Manual struct {
		AssistantRoleName string   `mapstructure:"assistant_role_name"`
		UserRoleName      string   `mapstructure:"user_role_name"`
		PhasePrompt       []string `mapstructure:"phase_prompt"`
	} `mapstructure:"Manual"`
}

func NewConfig() *Config {
	applicationConfig := &ApplicationConfig{}
	applicationConfig.readApplicationConfig()
	return &Config{
		ApplicationConfig: applicationConfig,
	}
}

func (c *ApplicationConfig) readApplicationConfig() {

	v := viper.New()
	v.SetConfigType("json")
	v.AddConfigPath("./configs/")

	// Read ChatChainConfig
	v.SetConfigName("ChatChainConfig")
	if err := v.ReadInConfig(); err != nil {
		utils.Logger.Panic().Err(err).Msg("Error reading ChatChainConfig")
	}
	if err := v.Unmarshal(&c.ChatDev.ChatChainConfig); err != nil {
		utils.Logger.Panic().Err(err).Msg("Error unmarshalling ChatChainConfig")
	}

	// Read PhaseConfig
	v.SetConfigName("PhaseConfig")
	if err := v.ReadInConfig(); err != nil {
		utils.Logger.Panic().Err(err).Msg("Error reading PhaseConfig")
	}
	if err := v.Unmarshal(&c.ChatDev.PhaseConfig); err != nil {
		utils.Logger.Panic().Err(err).Msg("Error unmarshalling PhaseConfig")
	}

	// Read RoleConfig
	v.SetConfigName("RoleConfig")
	if err := v.ReadInConfig(); err != nil {
		utils.Logger.Panic().Err(err).Msg("Error reading RoleConfig")
	}
	if err := v.Unmarshal(&c.ChatDev.RoleConfig); err != nil {
		utils.Logger.Panic().Err(err).Msg("Error unmarshalling RoleConfig")
	}
}
