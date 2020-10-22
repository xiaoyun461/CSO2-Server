package shop

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnShopList(p *PacketData, client net.Conn) {
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to request shop list but not in server !")
		return
	}
	//发送数据
	rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeShop), BuildShopList())
	SendPacket(rst, uPtr.CurrentConnection)
	DebugInfo(2, "Send a null shop list to User", string(uPtr.UserName))

}
func BuildShopList() []byte {
	return []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
}
