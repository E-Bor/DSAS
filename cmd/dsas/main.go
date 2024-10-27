package main

import (
	"DSAS/configs"
	"DSAS/internal/app"
	"DSAS/internal/config"
	"DSAS/pkg/config_loader"
	"DSAS/pkg/dsas_logger"
	"log/slog"
)

func main() {
	mainConfig, dsasConfig, err := mustLoad()

	if err != nil {
		return
	}
	dsas, err := app.NewApp(
		slog.Default(),
		mainConfig,
		dsasConfig,
	)

	err = dsas.Start()
	if err != nil {
		slog.Error(err.Error())
	}
}

func mustLoad() (
	*configs.MainConfig,
	*config.DsasConfig,
	error,
) {
	mainConfig, err := config_loader.LoadConfig(
		"configs/config.yaml",
		&configs.MainConfig{},
	)
	if err != nil {
		slog.Error(err.Error())
		return nil, nil, err
	}

	dsasConfig, err := config_loader.LoadConfig(
		mainConfig.ConfigPath.DsasConfig,
		&config.DsasConfig{},
	)

	if err != nil {
		slog.Error(err.Error())
		return nil, nil, err
	}

	dsas_logger.SetDefaultSlog(mainConfig.Env)

	return mainConfig, dsasConfig, nil
}
