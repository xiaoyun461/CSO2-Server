package chat

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnChatTeamMessage(p *InChatPacket, client net.Conn) {
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent TeamMessage but not in server !")
		return
	}
	if uPtr.GetUserRoomID() <= 0 {
		DebugInfo(2, "Error : User", string(uPtr.IngameName), "sent TeamMessage but not in room !")
		return
	}
	//找到对应房间
	rm := GetRoomFromID(uPtr.GetUserChannelServerID(), uPtr.GetUserChannelID(), uPtr.GetUserRoomID())
	if rm == nil || rm.Id <= 0 {
		DebugInfo(2, "Error : User", string(uPtr.IngameName), "sent TeamMessage but not in room !")
		return
	}
	if !uPtr.CurrentIsIngame {
		DebugInfo(2, "Error : User", string(uPtr.IngameName), "sent TeamMessage but not ingame !")
		return
	}
	//发送数据
	msg := BuildTeamMessage(uPtr, p)
	for _, v := range rm.Users {
		if v.CurrentIsIngame && v.GetUserTeam() == uPtr.GetUserTeam() {
			SendPacket(BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeChat), msg), v.CurrentConnection)
		}
	}
	DebugInfo(1, "User", string(uPtr.IngameName), "say team message <", string(p.Message), "> in room id", rm.Id)
}

func BuildTeamMessage(sender *User, p *InChatPacket) []byte {
	temp := make([]byte, 10+len(sender.IngameName)+int(p.MessageLen))
	offset := 0
	WriteUint8(&temp, ChatIngameTeam, &offset)
	WriteUint8(&temp, sender.Gm, &offset)
	WriteString(&temp, sender.IngameName, &offset)

	if sender.IsVIP() {
		WriteUint8(&temp, 1, &offset)
	} else {
		WriteUint8(&temp, 0, &offset)
	}
	WriteUint8(&temp, sender.VipLevel, &offset)

	WriteLongString(&temp, p.Message, &offset)
	return temp[:offset]

}
