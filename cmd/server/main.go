package main

import (
	"context"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/SuperJinggg/ai-router/internal/adapter"
	"github.com/SuperJinggg/ai-router/internal/config"
	"github.com/SuperJinggg/ai-router/internal/controller"
	"github.com/SuperJinggg/ai-router/internal/model/entity"
	"github.com/SuperJinggg/ai-router/internal/repository"
	"github.com/SuperJinggg/ai-router/internal/router"
	"github.com/SuperJinggg/ai-router/internal/service"
	"github.com/SuperJinggg/ai-router/internal/strategy"
	"github.com/SuperJinggg/ai-router/internal/task"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("open mysql with gorm failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("get sql db failed: %v", err)
	}
	defer sqlDB.Close()

	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("ping mysql failed: %v", err)
	}

	if err = autoMigrate(db); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}
	if err = seedData(db); err != nil {
		log.Fatalf("seed data failed: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	apiKeyRepo := repository.NewApiKeyRepository(db)
	requestLogRepo := repository.NewRequestLogRepository(db)
	providerRepo := repository.NewProviderRepository(db)
	modelRepo := repository.NewModelRepository(db)
	pluginRepo := repository.NewPluginRepository(db)
	userProviderKeyRepo := repository.NewUserProviderKeyRepository(db)
	rechargeRecordRepo := repository.NewRechargeRecordRepository(db)
	billingRecordRepo := repository.NewBillingRecordRepository(db)
	imageGenerationRecordRepo := repository.NewImageGenerationRecordRepository(db)
	redisPool := service.NewRedisPool(cfg)
	defer redisPool.Close()

	userService := service.NewUserService(userRepo)
	apiKeyService := service.NewApiKeyService(apiKeyRepo)
	requestLogService := service.NewRequestLogService(requestLogRepo, apiKeyService)
	billingService := service.NewBillingService(requestLogService)
	billingRecordService := service.NewBillingRecordService(billingRecordRepo)
	balanceService := service.NewBalanceService(userRepo, billingRecordService)
	rechargeService := service.NewRechargeService(rechargeRecordRepo, balanceService)
	stripePaymentService := service.NewStripePaymentService(cfg, rechargeService)
	chatCacheService := service.NewChatCacheService(redisPool, cfg)
	providerService := service.NewProviderService(providerRepo)
	userProviderKeyService := service.NewUserProviderKeyService(userProviderKeyRepo, providerService, cfg)
	pluginService := service.NewPluginService(pluginRepo, providerService, cfg)
	if err = pluginService.InitPlugins(); err != nil {
		log.Fatalf("init plugins failed: %v", err)
	}
	modelService := service.NewModelService(modelRepo, providerRepo)
	imageGenerationService := service.NewImageGenerationService(imageGenerationRecordRepo, modelService, providerService, userService, balanceService)
	healthCheckService := service.NewHealthCheckService(providerRepo, modelRepo, requestLogRepo)
	blacklistService := service.NewBlacklistService(redisPool)
	rateLimitService := service.NewRateLimitService(redisPool)

	routingStrategies := []strategy.RoutingStrategy{
		strategy.NewAutoRoutingStrategy(),
		strategy.NewFixedRoutingStrategy(),
		strategy.NewCostFirstRoutingStrategy(),
		strategy.NewLatencyFirstRoutingStrategy(),
	}
	routingService := service.NewRoutingService(modelRepo, routingStrategies)
	adapterFactory := adapter.NewModelAdapterFactory(
		[]adapter.ModelAdapter{
			adapter.NewZhipuAdapter(),
			adapter.NewOpenAIAdapter(),
		},
		adapter.NewDefaultAdapter(),
	)
	modelInvokeService := service.NewModelInvokeService(adapterFactory)
	chatService := service.NewChatService(requestLogService, routingService, modelInvokeService, providerService, userProviderKeyService, pluginService, userService, balanceService, chatCacheService)

	healthController := controller.NewHealthController()
	userController := controller.NewUserController(userService, requestLogService, billingService)
	apiKeyController := controller.NewApiKeyController(apiKeyService, userService)
	providerController := controller.NewProviderController(providerService)
	modelController := controller.NewModelController(modelService)
	blacklistController := controller.NewBlacklistController(blacklistService)
	chatController := controller.NewChatController(chatService, apiKeyService)
	internalChatController := controller.NewInternalChatController(chatService, userService)
	statsController := controller.NewStatsController(requestLogService, userService, billingService)
	pluginController := controller.NewPluginController(pluginService, userService)
	userProviderKeyController := controller.NewUserProviderKeyController(userProviderKeyService, userService)
	rechargeController := controller.NewRechargeController(rechargeService, stripePaymentService, userService, cfg)
	stripeWebhookController := controller.NewStripeWebhookController(stripePaymentService)
	balanceController := controller.NewBalanceController(balanceService, billingRecordService, userService)
	imageController := controller.NewImageController(imageGenerationService, userService, apiKeyService)
	healthCheckTask := task.NewHealthCheckTask(healthCheckService)

	engine, err := router.New(
		cfg,
		healthController,
		userController,
		apiKeyController,
		providerController,
		modelController,
		blacklistController,
		chatController,
		internalChatController,
		statsController,
		pluginController,
		userProviderKeyController,
		rechargeController,
		stripeWebhookController,
		balanceController,
		imageController,
		userService,
		blacklistService,
		rateLimitService,
	)
	if err != nil {
		log.Fatalf("build router failed: %v", err)
	}

	taskCtx, cancelTask := context.WithCancel(context.Background())
	defer cancelTask()
	healthCheckTask.Start(taskCtx)

	if err = engine.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server run failed: %v", err)
	}
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.ApiKey{},
		&entity.Model{},
		&entity.ModelProvider{},
		&entity.PluginConfig{},
		&entity.RequestLog{},
		&entity.BillingRecord{},
		&entity.RechargeRecord{},
		&entity.ImageGenerationRecord{},
		&entity.UserProviderKey{},
	)
}

func seedData(db *gorm.DB) error {
	var count int64
	db.Model(&entity.User{}).Count(&count)
	if count > 0 {
		return nil
	}

	users := []entity.User{
		{
			UserAccount:  "admin",
			UserPassword: "10670d38ec32fa8102be6a37f8cb52bf",
			UserName:     "管理员",
			UserAvatar:   "https://avatars.githubusercontent.com/u/58348976?v=4&size=64",
			UserProfile:  "系统管理员",
			UserRole:     "admin",
			UserStatus:   "active",
		},
		{
			UserAccount:  "user",
			UserPassword: "10670d38ec32fa8102be6a37f8cb52bf",
			UserName:     "普通用户",
			UserAvatar:   "https://avatars.githubusercontent.com/u/58348976?v=4&size=64",
			UserProfile:  "我是一个普通用户",
			UserRole:     "user",
			UserStatus:   "active",
		},
	}

	return db.Select(
		"userAccount", "userPassword", "userName", "userAvatar", "userProfile", "userRole", "userStatus",
	).Create(&users).Error
}
