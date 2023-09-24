package usecase

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"go-jobpass-bot/srcs/entities"
	"strconv"

	"go-jobpass-bot/srcs/store"
	"go-jobpass-bot/srcs/tools"
	"strings"
)

type Command struct {
	Command string
	Args    []string
	session *discordgo.Session
	message *discordgo.MessageCreate
}

var (
	CommandList = []discordgo.ApplicationCommand{
		{
			Name:        "register-gitlab",
			Description: "Link your discord account to your gitlab account",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        3,
					Name:        "username",
					Description: "Your gitlab username",
					Required:    true,
				},
			},
		},
		{
			Name:        "unlink-gitlab",
			Description: "Unlink your gitlab account from discord account",
		},
		{
			Name:        "list-gitlab-users",
			Description: "List all gitlab users related to project",
		},
		{
			Name:        "delete-discord-commands",
			Description: "Remove all discord Jobpass commands from Discord Server",
		},
		{
			Name:        "refresh-discord-commands",
			Description: "Refresh all discord Jobpass commands from Discord Server",
		},
	}
	CommandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"register-gitlab": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				message, err := RegisterGitlabAccount(s, i.ApplicationCommandData().Options[0].StringValue(), i.Member.User.ID)
				response := discordgo.InteractionResponse{
					Type: 4,
					Data: &discordgo.InteractionResponseData{
						Content: message,
						Title:   "register-gitlab-response",
					},
				}
				err = s.InteractionRespond(i.Interaction, &response)
				if err != nil {
					return
				}
				err = store.SaveStoreInfo()
				if err != nil {
					return
				}
			}
		},
		"unlink-gitlab": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				message, err := UnlinkGitlabAccount(s, i.Member.User.ID)
				response := discordgo.InteractionResponse{
					Type: 4,
					Data: &discordgo.InteractionResponseData{
						Content: message,
						Title:   "register-gitlab-response",
					},
				}
				_ = s.InteractionRespond(i.Interaction, &response)
				if err != nil {
					return
				}
				err = store.SaveStoreInfo()
				if err != nil {
					return
				}
			}
		},
		"list-gitlab-users": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				message, err := ListGitlabUsers()
				response := discordgo.InteractionResponse{
					Type: 4,
					Data: &discordgo.InteractionResponseData{
						Content: message,
						Title:   "register-gitlab-response",
					},
				}
				_ = s.InteractionRespond(i.Interaction, &response)
				if err != nil {
					return
				}
			}
		},
		"delete-discord-commands": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			errString := "Deleted all application commands"
			err := DeleteStoreCommands(s)
			if err != nil {
				errString = err.Error()
			}
			response := discordgo.InteractionResponse{
				Type: 4,
				Data: &discordgo.InteractionResponseData{
					Content: errString,
					Title:   "delete-command-response",
				},
			}
			err = s.InteractionRespond(i.Interaction, &response)
			if err != nil {
				return
			}
		},
		"refresh-discord-commands": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			errString := "Refreshed all application commands"
			err := RefreshStoreCommands(s)
			if err != nil {
				errString = err.Error()
			}
			response := discordgo.InteractionResponse{
				Type: 4,
				Data: &discordgo.InteractionResponseData{
					Content: errString,
					Title:   "refresh-command-response",
				},
			}
			err = s.InteractionRespond(i.Interaction, &response)
			if err != nil {
				return
			}
		},
	}
)

func InitCommands(s *discordgo.Session) error {
	if entities.Data.CommandInitialized == true {
		return nil
	}
	log.Infof("Creating discord commands")
	entities.Data.CommandInitialized = true
	for _, cmd := range CommandList {
		newCmd, err := s.ApplicationCommandCreate(entities.Data.Discord.ApplicationID, entities.Data.Discord.GuildID, &cmd)
		if err != nil {
			log.Fatalf("Cannot create slash command %q: %v", cmd.Name, err)
		} else {
			log.Infof("Created command %s", cmd.Name)
		}
		entities.Data.Commands = append(entities.Data.Commands, &entities.Command{
			ID:          newCmd.ID,
			Name:        newCmd.Name,
			Description: newCmd.Description,
		})
	}

	err := store.SaveStoreInfo()
	if err != nil {
		log.Infof("Error while saving store info")
		return err
	}
	return nil
}

func DeleteStoreCommands(s *discordgo.Session) error {
	for _, cmd := range entities.Data.Commands {
		err := s.ApplicationCommandDelete(entities.Data.Discord.ApplicationID, entities.Data.Discord.GuildID, cmd.ID)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", cmd.Name, err)
		} else {
			log.Infof("Deleted command %q", cmd.Name)
		}
	}
	entities.Data.Commands = nil
	entities.Data.CommandInitialized = false
	err := store.SaveStoreInfo()
	if err != nil {
		log.Infof("Error while saving store info")
		return err
	}
	return nil
}

func RefreshStoreCommands(s *discordgo.Session) error {
	err := DeleteStoreCommands(s)
	if err != nil {
		return err
	}
	err = InitCommands(s)
	if err != nil {
		log.Errorf("Error while loading commands | %s", err)
		return err
	}
	return nil
}

func UnlinkGitlabAccount(s *discordgo.Session, discordUserID string) (string, error) {
	discordUser, err := s.User(discordUserID)
	if err != nil {
		return "", err
	}

	boxTitle := fmt.Sprintf("Unlinking account")
	finalString := ""

	for i, user := range entities.Data.Users {
		if user.DiscordID == discordUserID {
			entities.Data.Users = append(entities.Data.Users[:i], entities.Data.Users[i+1:]...)
			finalString = fmt.Sprintf("Discord user %s is not linked to gitlab %s anymore", discordUser.Username, user.GitlabUsername)
			return tools.CreateTextBox(finalString, boxTitle), nil
		}
	}
	return "", nil
}

func RegisterGitlabAccount(s *discordgo.Session, gitlabUsername string, discordUserID string) (string, error) {
	boxTitle := fmt.Sprintf("Linking discord to gitlab user '%s'", gitlabUsername)
	finalString := ""

	for _, user := range entities.Data.Users {
		if user.DiscordID == discordUserID {
			boxTitle = "Error"
			finalString = fmt.Sprintf("Discord user %s is already linked to gitlab user %s", user.DiscordUsername, user.GitlabUsername)
			return tools.CreateTextBox(finalString, boxTitle), nil
		}
	}

	err := FetchGitlabProjectMembers()
	if err != nil {
		log.Errorf("Couldn't fetch gitlab users | %s", err)
		boxTitle = "Error"
		finalString += fmt.Sprintf("Internal system error")
		return tools.CreateTextBox(finalString, boxTitle), nil
	}
	gitlabUsername = strings.TrimPrefix(gitlabUsername, "@")

	for _, member := range entities.GitlabProjectMember {
		if strings.Compare(member.Username, gitlabUsername) == 0 {
			for _, storedUsers := range entities.Data.Users {
				if storedUsers.GitlabID == member.ID {
					discordUser, err := s.User(storedUsers.DiscordID)
					if err != nil {
						return "", err
					}
					boxTitle = "Error"
					finalString += fmt.Sprintf("Gitlab %s is already linked to discord %s\n", gitlabUsername, discordUser.Username)
					return tools.CreateTextBox(finalString, boxTitle), nil
				}
			}

			discordUser, err := s.User(discordUserID)
			if err != nil {
				boxTitle = "Error"
				finalString += fmt.Sprintf("%s\n", err)
				return tools.CreateTextBox(finalString, boxTitle), nil
			}
			log.Infof("Linking Gitlab account %d to discord account %s", member.ID, discordUserID)
			newUser := entities.User{
				ID:              strconv.Itoa(len(entities.Data.Users) + 1),
				GitlabID:        member.ID,
				GitlabUsername:  gitlabUsername,
				DiscordID:       discordUserID,
				DiscordUsername: discordUser.Username,
			}
			entities.Data.Users = append(entities.Data.Users, newUser)
			finalString += fmt.Sprintf(" - Discord:\n")
			finalString += fmt.Sprintf("   - ID: %s\n", newUser.DiscordID)
			finalString += fmt.Sprintf("   - Username: %s\n", newUser.DiscordUsername)
			finalString += fmt.Sprintf(" - Gitlab:\n")
			finalString += fmt.Sprintf("   - ID: %d\n", newUser.GitlabID)
			finalString += fmt.Sprintf("   - Username: %s\n", newUser.GitlabUsername)

			return tools.CreateTextBox(finalString, boxTitle), nil
		}
	}
	boxTitle = "Error"

	finalString += fmt.Sprintf("No gitlab user found with username %s\n", gitlabUsername)

	return tools.CreateTextBox(finalString, boxTitle), nil
}

func ListGitlabUsers() (string, error) {
	err := FetchGitlabProjectMembers()
	if err != nil {
		return "", err
	}
	finalString := ""
	userNames := make([]string, 0)
	longest := 0
	for _, gitUser := range entities.GitlabProjectMember {
		userNames = append(userNames, gitUser.Username)
		if len(gitUser.Username) > longest {
			longest = len(gitUser.Username)
		}
	}
	for _, name := range userNames {
		isLinked := false
		for _, user := range entities.Data.Users {
			if user.GitlabUsername == name {
				isLinked = true
			}
		}
		if isLinked == true {
			finalString += fmt.Sprintf("%s%*s | linked\n", name, longest-len(name), "")
		} else {
			finalString += fmt.Sprintf("%s%*s |       \n", name, longest-len(name), "")
		}
	}

	return tools.CreateTextBox(finalString, "Gitlab Users"), nil
}
