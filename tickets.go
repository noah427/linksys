package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

type ticket struct {
	userID          string
	guildID         string
	ticketChannelID string
	cost            int
	info            string
	ticketCloserID  string
	closingProgress int
	creationTime    time.Time
}

func (t *ticket) getUsername() string {
	user, err := client.User(t.userID)

	if err != nil {
		log.Println("User Not Found")
		return "User Not Found"
	}

	return fmt.Sprintf("%s#%s", user.Username, user.Discriminator)
}

func (t *ticket) getGuildName() string {
	guild, err := client.Guild(t.guildID)
	if err != nil {
		log.Println("Guild Not Found")
		return "Guild Not Found"
	}
	return guild.Name
}

func (t *ticket) generateEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Linksys | Sale",
		Color: 238,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Client", Value: t.getUsername()},
			{Name: "Selling Server", Value: t.getGuildName()},
			{Name: "Price Paid", Value: fmt.Sprintf("%s$", strconv.Itoa(t.cost))},
			{Name: "Amount Owed to network", Value: fmt.Sprintf("%s$", strconv.Itoa(t.cost/2))},
			{Name: "Info", Value: t.info},
		},
	}

}

func (t *ticket) generateIncompleteEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: 238,
		Title: fmt.Sprintf("%s's ticket", t.getUsername()),
	}

}

func (t *ticket) deleteSelf() {
	for i, tic := range tickets {
		if tic.guildID == t.guildID && tic.ticketChannelID == t.ticketChannelID {
			tickets = remove(tickets, i)
		}
	}
}
