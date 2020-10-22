package shop

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

const (
	requestList = 3
)

func OnShopRequest(p *PacketData, client net.Conn) {
	var pkt InShopPacket
	if p.PraseShopPacket(&pkt) {
		switch pkt.InShopType {
		case requestList:
			OnShopList(p, client)
		default:
			DebugInfo(2, "Unknown shop packet", pkt.InShopType, "from", client.RemoteAddr().String(), p.Data)
		}
	} else {
		DebugInfo(2, "Error : Recived a illegal shop packet from", client.RemoteAddr().String())
	}
}
