package strategy

import "github.com/SuperJinggg/ai-router/internal/model/entity"

type RoutingStrategy interface {
	SelectModel(models []entity.Model, requestedModel string) *entity.Model
	GetFallbackModels(models []entity.Model, requestedModel string) []entity.Model
	GetStrategyType() string
}
