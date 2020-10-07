package achievement

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnAchievementCampaign(p *PacketData, client net.Conn) {
	//检索数据包
	var pkt InAchievementCampaignPacket
	if !p.PraseInAchievementCampaignPacket(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a error AchievementCampaign packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent AchievementCampaign packet but not in server !")
		return
	}
	//处理关卡
	switch pkt.CampaignId {
	case Campaign_0:
		reward := OutAchievementCampaign{}
		rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeAchievement), BuildAchievementCampaign(0, &reward, pkt.CampaignId))
		SendPacket(rst, uPtr.CurrentConnection)
		DebugInfo(1, "User", string(uPtr.UserName), "get achievement campaign id", pkt.CampaignId)
	case Campaign_1:
		reward := OutAchievementCampaign{0, 0, 0, 0, 0, 3000, 0, []OutAchievementCampaignItems{}, 0}
		rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeAchievement), BuildAchievementCampaign(0x20, &reward, pkt.CampaignId))
		SendPacket(rst, uPtr.CurrentConnection)
		DebugInfo(1, "User", string(uPtr.UserName), "get achievement campaign id", pkt.CampaignId)
	case Campaign_2:
		reward := OutAchievementCampaign{0, 0, 0, 0, 5000, 0, 0, []OutAchievementCampaignItems{}, 0}
		rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeAchievement), BuildAchievementCampaign(0x10, &reward, pkt.CampaignId))
		SendPacket(rst, uPtr.CurrentConnection)
		DebugInfo(1, "User", string(uPtr.UserName), "get achievement campaign id", pkt.CampaignId)
	case Campaign_3:
		reward := OutAchievementCampaign{0, 0, 0, 24, 0, 0, 0, []OutAchievementCampaignItems{}, 0}
		rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeAchievement), BuildAchievementCampaign(0x8, &reward, pkt.CampaignId))
		SendPacket(rst, uPtr.CurrentConnection)
		DebugInfo(1, "User", string(uPtr.UserName), "get achievement campaign id", pkt.CampaignId)
	case Campaign_4:
		reward := OutAchievementCampaign{0, 0, 0, 0, 0, 0, 2, []OutAchievementCampaignItems{{1002, 1, 0}, {1004, 1, 0}}, 0}
		rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeAchievement), BuildAchievementCampaign(0x40, &reward, pkt.CampaignId))
		SendPacket(rst, uPtr.CurrentConnection)
		DebugInfo(1, "User", string(uPtr.UserName), "get achievement campaign id", pkt.CampaignId)
	case Campaign_5:
		reward := OutAchievementCampaign{0, 0, 0, 0, 0, 0, 2, []OutAchievementCampaignItems{{54, 1, 0}, {55, 1, 0}}, 0}
		rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeAchievement), BuildAchievementCampaign(0x40, &reward, pkt.CampaignId))
		SendPacket(rst, uPtr.CurrentConnection)
		DebugInfo(1, "User", string(uPtr.UserName), "get achievement campaign id", pkt.CampaignId)
	case Campaign_6:

		DebugInfo(1, "User", string(uPtr.UserName), "get achievement campaign id", pkt.CampaignId)
	default:
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a unkown AchievementCampaign packet id", pkt.CampaignId)
	}
}

func BuildAchievementCampaign(flags uint32, src *OutAchievementCampaign, id uint16) []byte {
	if src == nil {
		return []byte{}
	}
	buf := make([]byte, 25+src.NumOfItems*8)
	offset := 0
	WriteUint8(&buf, campaign, &offset)
	WriteUint16(&buf, id, &offset)
	WriteUint32(&buf, flags, &offset)

	if flags&0x1 != 0 {
		WriteUint16(&buf, src.Unk00, &offset)
	}
	if flags&0x2 != 0 {
		WriteUint32(&buf, src.Unk01, &offset)
	}
	if flags&0x4 != 0 {
		WriteUint16(&buf, src.RewardTitle, &offset)
	}
	if flags&0x8 != 0 {
		WriteUint16(&buf, src.RewardIcon, &offset)
	}
	if flags&0x10 != 0 {
		WriteUint32(&buf, src.RewardPoints, &offset)
	}
	if flags&0x20 != 0 {
		WriteUint32(&buf, src.RewardXp, &offset)
	}
	if flags&0x40 != 0 {
		WriteUint8(&buf, src.NumOfItems, &offset)
		for _, v := range src.Items {
			WriteUint32(&buf, v.ItemId, &offset)
			WriteUint16(&buf, v.Ammount, &offset)
			WriteUint16(&buf, v.TimeLimited, &offset)
		}
	}
	if flags&0x80 != 0 {
		WriteUint16(&buf, src.Unk02, &offset)
	}
	return buf[:offset]
}
