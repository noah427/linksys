package main

import (
	"errors"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var specialChars = regexp.MustCompile(`([^a-z \d])`)

var errInsufficientPerms = errors.New("Insufficient Perms")

func getEveryoneRoleID(roles []*discordgo.Role) string {
	for _, role := range(roles){
		if role.Name == "@everyone"{
			return role.ID
		}
	}
	return ""
}

func findTicketCategory(s *discordgo.Session, guildID string) (channelID string, err error) {
	channels, _ := s.GuildChannels(guildID)

	for _, channel := range channels {
		if channel.Name == "tickets" && channel.Type == discordgo.ChannelTypeGuildCategory {
			return channel.ID, nil
		}
	}

	guild, _ := s.Guild(guildID)


	channel, err := s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Type : discordgo.ChannelTypeGuildCategory,
		Name : "tickets",
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{ID : getEveryoneRoleID(guild.Roles), Deny : 1024, Type : "role"},
		},
	})

	if err != nil {
		return "", errInsufficientPerms
	}

	return channel.ID, nil
}

func getTicketByUsername(username string) *ticket {
	for _, tic := range tickets {
		comparable := strings.ToLower(tic.getUsername())
		comparable = removeSpecialChars(comparable)

		if comparable == removeSpecialChars(username) {
			return tic
		}
	}

	return nil
}

func removeSpecialChars(username string) string {
	return specialChars.ReplaceAllString(strings.ToLower(username), "")
}


func remove(s []*ticket, i int) []*ticket {
    s[len(s)-1], s[i] = s[i], s[len(s)-1]
    return s[:len(s)-1]
}