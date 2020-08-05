package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func messageCreate(s *discordgo.Session, msg *discordgo.MessageCreate) {
	prefix := os.Getenv("PREFIX")
	if msg.Author.ID == s.State.User.ID {
		return
	}

	detectClosingTicket(msg.Message, s)

	if strings.HasPrefix(msg.Content, prefix) {
		command := strings.TrimPrefix(msg.Content, prefix)
		command = strings.ToLower(command)

		if command == "help" || command == "info" {
			s.ChannelMessageSendEmbed(msg.ChannelID, helpEmbed)
		}

		if command == "close" {
			channel, _ := s.Channel(msg.ChannelID)

			username := strings.Replace(channel.Name, "-", "#", 1)
			tic := getTicketByUsername(username)

			// if username == fmt.Sprintf("%s#%s", removeSpecialChars(msg.Author.Username), msg.Author.Discriminator) {
			// 	s.ChannelMessageSend(msg.ChannelID, "You may not close this ticket")
			// 	return
			// }

			if tic == nil {
				s.ChannelMessageSend(msg.ChannelID, "This is not a ticket channel")
				return
			}

			tic.ticketCloserID = msg.Author.ID

			tic.closeTicket(s, nil)
		}

		if command == "ticket" {

			ticketCatID, err := findTicketCategory(s, msg.GuildID)

			if err != nil {
				s.ChannelMessageSend(msg.ChannelID, "Insufficient Permissions, please contact guild owner")
				return
			}

			user, _ := s.User(msg.Author.ID)
			username := fmt.Sprintf("%s#%s", user.Username, user.Discriminator)
			if getTicketByUsername(username) != nil {
				m, _ := s.ChannelMessageSend(msg.ChannelID, "You may not have multiple tickets open")
				time.AfterFunc(time.Second*2, func() { s.ChannelMessageDelete(msg.ChannelID, m.ID) })
				return
			}

			m, _ := s.ChannelMessageSend(msg.ChannelID, "Creating ticket...")
			time.AfterFunc(time.Second*2, func() { s.ChannelMessageDelete(msg.ChannelID, m.ID) })

			channel, _ := s.GuildChannelCreateComplex(msg.GuildID, discordgo.GuildChannelCreateData{
				Name:     fmt.Sprintf("%s-%s", msg.Author.Username, msg.Author.Discriminator),
				Type:     discordgo.ChannelTypeGuildText,
				ParentID: ticketCatID,
				PermissionOverwrites: []*discordgo.PermissionOverwrite{
					{ID: msg.Author.ID, Allow: 1024, Type: "member"},
				},
			})

			tic := &ticket{
				userID:          msg.Author.ID,
				guildID:         msg.GuildID,
				ticketChannelID: channel.ID,
				creationTime:    time.Now(),
			}

			s.ChannelMessageSendEmbed(channel.ID, tic.generateIncompleteEmbed())

			tickets = append(tickets, tic)
		}

	}

}
