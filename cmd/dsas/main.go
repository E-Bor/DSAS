package main

import (
	"DSAS2/configs"
	"DSAS2/internal/app/reports_registry"
	"DSAS2/internal/config"
	"DSAS2/pkg/config_loader"
	"DSAS2/pkg/dsas_logger"
	"fmt"
	"log/slog"
)

func main() {
	dsas_logger.SetDefaultSlog()

	mainConf, err := config_loader.LoadConfig(
		"configs/config.yaml",
		&configs.MainConfig{},
	)
	if err != nil {
		return
	}

	dsasConfig, err := config_loader.LoadConfig(
		mainConf.ConfigPath.DsasConfig,
		&config.DsasConfig{},
	)

	if err != nil {
		return
	}

	reportMap, err := reports_registry.NewReportRegistry(dsasConfig.IntegrationsDir)

	fmt.Println(reportMap)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	rep1, err := reportMap.Get(
		"facebook",
		"get_pet_report",
	)
	rep3, err := reportMap.Get(
		"facebook",
		"get_store_report",
	)

	rep2, err := reportMap.Get(
		"instagram",
		"get_store_report",
	)

	err = rep2()
	err = rep1()
	err = rep3()

	if err != nil {
		return
	}
}
