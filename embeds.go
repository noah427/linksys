package main

import (
	"github.com/bwmarrin/discordgo"
)

var helpEmbed = &discordgo.MessageEmbed{
	Title: "Linksys | Help",
	Color: 238,
	Fields: []*discordgo.MessageEmbedField{
		{Name : "help : ", Value: "shows this message", Inline : true},
		{Name : "ticket : ", Value: "creates a ticket", Inline: true},
		{Name : "close : ", Value: "closes a ticket (must be done in a ticket channel) (cannot be done by the creator of the ticket)"},
	},
}

