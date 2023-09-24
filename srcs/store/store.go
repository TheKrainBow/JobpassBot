package store

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
	"go-jobpass-bot/srcs/entities"
	"go-jobpass-bot/srcs/tools"
)

func InitStore() error {
	log.Info("Recovering data ...")
	dat, _ := os.ReadFile(entities.StorePath)
	err := json.Unmarshal(dat, &entities.Data)
	if err != nil {
		return err
	}
	tools.LogDeleteLastNLines(1)
	log.Info("Recovering data ✓")
	return nil
}

func SaveStoreInfo() error {
	log.Info("Saving data ...")
	entities.DataMutext.Lock()
	defer func() {
		entities.DataMutext.Unlock()
	}()
	save := entities.Data

	data, err := json.MarshalIndent(save, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(entities.StorePath, data, 0o644)
	if err != nil {
		return err
	}
	tools.LogDeleteLastNLines(1)
	log.Info("Saving data ✓")
	return nil
}
