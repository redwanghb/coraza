package types

import "fmt"

type LLMAction int

const (
	// LLMActionBlock为阻断LLM请求或回复，直接断开连接
	LLMActionBlock LLMAction = iota + 1
	// LLMActionWarn为通过对话框发送告警通知，不回复LLM回复的内容
	LLMActionWarn
	// LLMActionDesentize为对LLM回复的内容进行脱敏处理
	LLMActionDesentize
	// LLMActionAllow为允许LLM请求或回复
	LLMActionAllow
)

func (l LLMAction) String() string {
	switch l {
	case LLMActionBlock:
		return "block"
	case LLMActionWarn:
		return "warn"
	case LLMActionDesentize:
		return "desentize"
	case LLMActionAllow:
		return "allow"
	default:
		return "unknown"
	}
}

type LLMCategoryConfig struct {
	// 对应规则ID
	ID int `json:"id" yaml:"id" export:"true"`
	// 对应分类名称
	Name string `json:"name" yaml:"name" export:"true"`
	// 对应规则的严重级别
	// 0: info, 1: low, 2: medium, 3: high， 4: critical
	Severity RuleSeverity `json:"severity" yaml:"severity" export:"true"`
	// 对应规则的动作
	// 1: block, 2: warn, 3: desentize, 4: allow， 5: unknown
	LLMAction LLMAction `json:"llmAction" yaml:"llmAction" export:"true"`
	// 对应规则的标签
	Tags []string `json:"tags" yaml:"tags" export:"true"`
}

func NewLLMCategoryConfig(id int, name string, severity RuleSeverity, action LLMAction, tags []string) *LLMCategoryConfig {
	return &LLMCategoryConfig{
		ID:        id,
		Name:      name,
		Severity:  severity,
		LLMAction: action,
		Tags:      tags,
	}
}

func (categroy *LLMCategoryConfig) GetAction() LLMAction {
	return categroy.LLMAction
}

func (category *LLMCategoryConfig) GetSeverity() RuleSeverity {
	return category.Severity
}

func (category *LLMCategoryConfig) GetTags() []string {
	return category.Tags
}

func (category *LLMCategoryConfig) GetID() int {
	return category.ID
}

type ClassificationConfig struct {
	Categories map[string]LLMCategoryConfig `json:"categories" yaml:"categories" export:"true"`
}

func NewClassificationConfig() *ClassificationConfig {
	return &ClassificationConfig{
		Categories: make(map[string]LLMCategoryConfig),
	}
}

func (c *ClassificationConfig) AddCategory(category *LLMCategoryConfig) {
	c.Categories[category.Name] = *category
}

func (c *ClassificationConfig) LLMCategoryConfig(name string) (*LLMCategoryConfig, error) {
	categoryConfig, ok := c.Categories[name]
	if !ok {
		return nil, fmt.Errorf("category %s not found", name)
	}
	return &categoryConfig, nil
}

type LLMTypeSeverity struct {
	// 提示词注入检测结果
	PromptInjection RuleSeverity
	// 检测问题中有代码片段
	CodeSippetsInput RuleSeverity
	// 检测回答中有代码片段
	CodeSippetsOutput RuleSeverity
	// 检测回答中有敏感词
	SensitiveOutput RuleSeverity
	// 检测问题中包含了政治敏感词
	PoliticalInput RuleSeverity
	// 检测回答中包含了政治敏感词
	PoliticalOutput RuleSeverity
	// 检测问题中包含了歧视性、仇恨性、暴力性等违法内容
	IllegalInput RuleSeverity
	// 检测回答中包含了歧视性、仇恨性、暴力性等违法内容
	IllegalOutput RuleSeverity
	// 检测问题中包含了恶意链接
	MaliciousInput RuleSeverity
	// 检测回答中包含了恶意链接
	MaliciousOutput RuleSeverity
	// 检测问题中包含了不合规的内容
	NoncompliantInput RuleSeverity
	// 检测回答中包含了不合规的内容
	NoncompliantOutput RuleSeverity
}
