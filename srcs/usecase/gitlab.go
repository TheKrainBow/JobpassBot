package usecase

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"go-jobpass-bot/srcs/dto"
	"go-jobpass-bot/srcs/entities"
)

func FetchGitlabProjectMembers() error {
	if len(entities.GitlabProjectMember) != 0 {
		return nil
	}
	git, err := gitlab.NewClient(entities.Data.Gitlab.APIKey)
	if err != nil {
		log.Fatalf("Failed to create gitlab client: %v", err)
		return err
	}

	members, _, err := git.ProjectMembers.ListProjectMembers(entities.Data.Gitlab.ProjectID, nil)
	if err != nil {
		log.Errorf("wtf happened | %s", err)
		return err
	}
	for _, member := range members {
		if strings.Compare(member.Name, "discordBot") == 0 {
			continue
		}
		entities.GitlabProjectMember = append(entities.GitlabProjectMember, entities.GitlabUser{
			ID:          member.ID,
			Name:        member.Name,
			Username:    member.Username,
			Mail:        member.Email,
			AccessLevel: dto.GitlabAccessLevelToEntities(member.AccessLevel),
		})
	}
	return nil
}
