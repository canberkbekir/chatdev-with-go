package utils

import (
	"chatdev-with-go/internal/models"
	"github.com/spf13/viper"
	"strings"
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
	Roles map[string][]string `mapstructure:",remain"`
}

type Role struct {
	RoleName       string
	RolePromptFull string
	RolePrompt     []string
}

type PhaseConfig struct {
	Phases map[string]models.PhaseDetails `mapstructure:"phases"`
}

type Phase struct {
	PhaseName string
	Phase     models.PhaseDetails
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
		Logger.Panic().Err(err).Msg("Error reading ChatChainConfig")
	}
	if err := v.Unmarshal(&c.ChatDev.ChatChainConfig); err != nil {
		Logger.Panic().Err(err).Msg("Error unmarshalling ChatChainConfig")
	}

	// Read PhaseConfig
	v.SetConfigName("PhaseConfig")
	if err := v.ReadInConfig(); err != nil {
		Logger.Panic().Err(err).Msg("Error reading PhaseConfig")
	}
	if err := v.Unmarshal(&c.ChatDev.PhaseConfig); err != nil {
		Logger.Panic().Err(err).Msg("Error unmarshalling PhaseConfig")
	}

	// Read RoleConfig
	v.SetConfigName("RoleConfig")
	if err := v.ReadInConfig(); err != nil {
		Logger.Panic().Err(err).Msg("Error reading RoleConfig")
	}
	if err := v.Unmarshal(&c.ChatDev.RoleConfig); err != nil {
		Logger.Panic().Err(err).Msg("Error unmarshalling RoleConfig")
	}
}

func NewRole(roleName string, rolePrompt []string) *Role {
	roleText := strings.Join(rolePrompt, ", ")
	return &Role{
		RoleName:       roleName,
		RolePromptFull: roleText,
		RolePrompt:     rolePrompt,
	}
}

func NewPhase(phaseName string, phase models.PhaseDetails) *Phase {
	return &Phase{
		PhaseName: phaseName,
		Phase:     phase,
	}
}
