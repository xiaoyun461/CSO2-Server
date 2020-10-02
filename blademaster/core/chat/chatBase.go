package chat

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnChat(p *PacketData, client net.Conn) {
	var pkt InChatPacket
	if p.PraseInChatPacket(&pkt) {
		switch pkt.Type {
		case ChatDirectMessage:
			OnChatDirectMessage(&pkt, client)
		case ChatChannel:
			OnChatChannelMessage(&pkt, client)
		case ChatRoom:

		case ChatIngameGlobal:

		case ChatIngameTeam:

		default:
			DebugInfo(2, "Unknown chat packet", pkt.Type, "from", client.RemoteAddr().String())
		}
	} else {
		DebugInfo(2, "Error : Recived a illegal host packet from", client.RemoteAddr().String())
	}
}
