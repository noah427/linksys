package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (t *ticket) closeTicket(s *discordgo.Session, response *string) {
	if response == nil && t.closingProgress != 0 {
		s.ChannelMessageSend(t.ticketChannelID, "Response may not be empty")
		return
	}

	switch t.closingProgress {
	case 0:
		s.ChannelMessageSend(t.ticketChannelID, "Please enter the price paid by the client (USD)")
		break
	case 1:
		cost, err := strconv.Atoi(*response)
		if err != nil {
			s.ChannelMessageSend(t.ticketChannelID, "Cost must be a number (USD)")
			return
		}
		t.cost = cost

		s.ChannelMessageSend(t.ticketChannelID, "Please supply extra information on what was purchased")
		break
	case 2:
		t.info = *response

		s.ChannelMessageSend(t.ticketChannelID, "Thank you! This channel will close in 10 seconds")

		time.AfterFunc(time.Second*10, func() {
			s.ChannelDelete(t.ticketChannelID)
			_, err := s.ChannelMessageSendEmbed(os.Getenv("LOGGINGCHANNELID"), t.generateEmbed())
			if err != nil {
				log.Println("Cannot send in logging channel")
			}
			t.deleteSelf()
		})

		break
	}

	t.closingProgress++
}

func detectClosingTicket(msg *discordgo.Message,s *discordgo.Session){
	for _,tic :=  range(tickets){
		if tic.ticketCloserID == msg.Author.ID && tic.ticketChannelID == msg.ChannelID {
			tic.closeTicket(s, &msg.Content)
		}
	}
}
