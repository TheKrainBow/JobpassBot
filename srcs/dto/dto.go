package dto

import (
	"github.com/xanzy/go-gitlab"
	"go-jobpass-bot/srcs/entities"
)

func GitlabAccessLevelToEntities(value gitlab.AccessLevelValue) entities.GitlabAccessLevel {
	switch value {
	case 0:
		return entities.NoPermissions
	case 5:
		return entities.MinimalAccessPermissions
	case 10:
		return entities.GuestPermissions
	case 20:
		return entities.ReporterPermissions
	case 30:
		return entities.DeveloperPermissions
	case 40:
		return entities.MaintainerPermissions
	case 50:
		return entities.OwnerPermissions
	case 60:
		return entities.AdminPermissions
	}
	return entities.NoPermissions
}

func EntitiesAccessLevelToGitlab(value entities.GitlabAccessLevel) gitlab.AccessLevelValue {
	switch value {
	case entities.NoPermissions:
		return 0
	case entities.MinimalAccessPermissions:
		return 5
	case entities.GuestPermissions:
		return 10
	case entities.ReporterPermissions:
		return 20
	case entities.DeveloperPermissions:
		return 30
	case entities.MaintainerPermissions:
		return 40
	case entities.OwnerPermissions:
		return 50
	case entities.AdminPermissions:
		return 60
	}
	return 0
}
