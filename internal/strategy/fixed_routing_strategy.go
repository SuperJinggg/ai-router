package strategy

import (
	"sort"

	"github.com/SuperJinggg/ai-router/internal/constant"
	"github.com/SuperJinggg/ai-router/internal/model/entity"
)

type FixedRoutingStrategy struct{}

func NewFixedRoutingStrategy() *FixedRoutingStrategy {
	return &FixedRoutingStrategy{}
}

func (s *FixedRoutingStrategy) GetStrategyType() string {
	return constant.RoutingStrategyFixed
}

func (s *FixedRoutingStrategy) SelectModel(models []entity.Model, requestedModel string) *entity.Model {
	if requestedModel == "" {
		return nil
	}
	for _, model := range models {
		if model.ModelKey == requestedModel {
			selected := model
			return &selected
		}
	}
	return nil
}

func (s *FixedRoutingStrategy) GetFallbackModels(models []entity.Model, requestedModel string) []entity.Model {
	cloned := append([]entity.Model(nil), models...)
	sort.Slice(cloned, func(i, j int) bool {
		if cloned[i].Priority == cloned[j].Priority {
			return cloned[i].ID < cloned[j].ID
		}
		return cloned[i].Priority > cloned[j].Priority
	})
	result := make([]entity.Model, 0, len(cloned))
	for _, model := range cloned {
		if requestedModel != "" && model.ModelKey == requestedModel {
			continue
		}
		result = append(result, model)
	}
	return result
}
