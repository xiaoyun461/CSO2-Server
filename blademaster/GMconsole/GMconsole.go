package GMconsole

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"

	. "github.com/KouKouChan/CSO2-Server/blademaster/core/message"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/room"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/configure"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

type GMInfo struct {
	GMip       string
	GMport     string
	GMname     string
	GMpassword string
}

var (
	clients map[string]bool
)

func InitGMconsole() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Console server suffered a fault !")
			fmt.Println("error:", err)
			fmt.Println("Fault end!")
		}
	}()
	clients = map[string]bool{}

	//init TCP
	GMserver, err := net.Listen("tcp", fmt.Sprintf(":%d", Conf.GMport))
	if err != nil {
		fmt.Println("Init console tcp socket error !\n")
		panic(err)
	}
	defer GMserver.Close()

	fmt.Println("Console is running at", "[AnyAdapter]:"+strconv.Itoa(int(Conf.GMport)))
	for {
		client, err := GMserver.Accept()
		if err != nil {
			DebugInfo(2, "Console server accept data error !\n")
			continue
		}
		DebugInfo(2, "Console server accept a new connection request at", client.RemoteAddr().String())
		go RecvGMmsg(client)
	}
}

func RecvGMmsg(client net.Conn) {
	defer client.Close()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("GM Client", client.RemoteAddr().String(), "suffered a fault !")
			fmt.Println(err)
			fmt.Println("Fault end!")
			delete(clients, client.RemoteAddr().String())
		}
	}()

	for {
		//读取3字节数据包头部
		headBytes, err := GMReadHead(client)
		if !err {
			goto end
		}
		if headBytes[0] != GMSignature {
			DebugInfo(2, "Recived a illegal GM head from", client.RemoteAddr().String())
			continue
		}
		offset := 1
		Len := ReadUint16(headBytes, &offset)
		//读取数据部分
		bytes, err := GMReadData(client, Len)
		if !err {
			goto end
		}
		dataPacket := GMpacket{
			bytes,
			Len,
			string(bytes),
		}
		req := strings.Fields(string(dataPacket.Req))

		switch req[0] {
		case GMLogin:
			login(client, req)
		case GMReqUserList:
			userlist(client, req)
		case GMKickUser:
			kickUser(client, req)
		case GMadditem:
			additem(client, req)
		case GMsave:
			save(client, req)
		default:
		}
	}
end:
	DebugInfo(1, "Console client", client.RemoteAddr().String(), "closed the connection")
	delete(clients, client.RemoteAddr().String())
	return
}

func login(client net.Conn, req []string) {
	if req[1] == Conf.GMusername && req[2] == Conf.GMpassword {
		clients[client.RemoteAddr().String()] = true
		rst := []byte(GMLoginOk)
		GMSendPacket(&rst, client)
		DebugInfo(1, "Console from", client.RemoteAddr().String(), "logged in")
	} else {
		rst := []byte(GMLoginFailed)
		GMSendPacket(&rst, client)
		DebugInfo(1, "Console from", client.RemoteAddr().String(), "login failed")
	}
}

func userlist(client net.Conn, req []string) {
	if _, ok := clients[client.RemoteAddr().String()]; !ok || !clients[client.RemoteAddr().String()] {
		DebugInfo(2, "Console from", client.RemoteAddr().String(), "sent a userlist req but not logged in")
	}
	rst := OutUserList{
		len(UsersManager.Users),
		[]string{},
	}
	for _, v := range UsersManager.Users {
		if v == nil {
			rst.UserNum--
			continue
		}
		rst.UserNames = append(rst.UserNames, string(v.UserName))
	}
	jsondata, _ := json.Marshal(rst)
	GMSendPacket(&jsondata, client)
	DebugInfo(1, "Console from", client.RemoteAddr().String(), "request a userlist")

}

func kickUser(client net.Conn, req []string) {
	if _, ok := clients[client.RemoteAddr().String()]; !ok || !clients[client.RemoteAddr().String()] {
		DebugInfo(2, "Console from", client.RemoteAddr().String(), "sent a kick user req but not logged in")
	}
	for _, v := range UsersManager.Users {
		if v == nil {
			continue
		}
		if string(v.UserName) == req[1] {
			OnSendMessage(v.CurrentSequence, v.CurrentConnection, MessageDialogBox, GAME_SERVER_ERROR)
			OnLeaveRoom(v.CurrentConnection, true)
			DelUserWithConn(v.CurrentConnection)
			v.CurrentConnection.Close()

			rst := []byte(GMKickSuccess)
			GMSendPacket(&rst, client)
			DebugInfo(1, "Console from", client.RemoteAddr().String(), "kicked player", string(v.UserName))
			return
		}
	}
	rst := []byte(GMKickFailed)
	GMSendPacket(&rst, client)
}

func additem(client net.Conn, req []string) {
	if _, ok := clients[client.RemoteAddr().String()]; !ok || !clients[client.RemoteAddr().String()] {
		DebugInfo(2, "Console from", client.RemoteAddr().String(), "sent a additem req but not logged in")
	}

	id, err := strconv.Atoi(req[2])
	if err != nil {
		rst := []byte(GMAdditemFailed)
		GMSendPacket(&rst, client)
	}

	for _, v := range UsersManager.Users {
		if v == nil {
			continue
		}
		if string(v.UserName) == req[1] {

			v.AddItem(uint32(id))

			rst := BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeUserInfo), BuildUserInfo(0XFFFFFFFF, NewUserInfo(v), v.Userid, true))
			SendPacket(rst, v.CurrentConnection)

			rst = []byte(GMAdditemSuccess)
			GMSendPacket(&rst, client)

			DebugInfo(1, "Console from", client.RemoteAddr().String(), "add item", id, "to User", string(v.UserName))
			return
		}
	}
	filepath := DBPath + req[1]
	rb, _ := PathExists(filepath)
	if rb {
		u := GetNewUser()

		Dblock.Lock()
		dataEncoded, _ := ioutil.ReadFile(filepath)
		Dblock.Unlock()

		err := json.Unmarshal(dataEncoded, &u)
		if err == nil {
			rst := []byte(GMAdditemSuccess)
			GMSendPacket(&rst, client)

			DebugInfo(1, "Console from", client.RemoteAddr().String(), "add item", id, "to User", string(u.UserName))
			return

		}

		u.AddItem(uint32(id))
		err = UpdateUserToDB(&u)

		if err == nil {
			rst := []byte(GMAdditemSuccess)
			GMSendPacket(&rst, client)

			DebugInfo(1, "Console from", client.RemoteAddr().String(), "add item", id, "to User", string(u.UserName))
			return

		}
	}
	rst := []byte(GMAdditemFailed)
	GMSendPacket(&rst, client)
}

func save(client net.Conn, req []string) {
	if _, ok := clients[client.RemoteAddr().String()]; !ok || !clients[client.RemoteAddr().String()] {
		DebugInfo(2, "Console from", client.RemoteAddr().String(), "sent a additem req but not logged in")
	}
	for _, v := range UsersManager.Users {
		if v == nil {
			continue
		}

		if UpdateUserToDB(v) != nil {
			rst := []byte(GMSaveFailed)
			GMSendPacket(&rst, client)
			return
		}
	}

	DebugInfo(1, "Console from", client.RemoteAddr().String(), "request to save all data")

	rst := []byte(GMSaveSuccess)
	GMSendPacket(&rst, client)
}
