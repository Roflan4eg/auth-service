package handlers

import "github.com/Roflan4eg/auth-serivce/config"

type Container struct {
	Metrics *MetricsHandler
	// Health *HealthHandler
}

func NewContainer(cfg *config.Config) *Container {
	return &Container{
		Metrics: NewMetricsHandler(cfg),
		// Health: NewHealthHandler(),
	}
}
