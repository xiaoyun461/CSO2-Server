package main

import (
	"flag"
	"fmt"
	"math"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	. "github.com/KouKouChan/CSO2-Server/blademaster/Exp"
	. "github.com/KouKouChan/CSO2-Server/blademaster/GMconsole"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/achievement"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/automatch"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/chat"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/holepunch"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/host"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/inventory"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/mail"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/message"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/option"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/playerinfo"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/pointlotto"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/quick"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/report"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/room"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/shop"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/supply"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/user"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/version"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/configure"
	. "github.com/KouKouChan/CSO2-Server/database/redis"
	. "github.com/KouKouChan/CSO2-Server/database/sqlite"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/kerlong/encode"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
	. "github.com/KouKouChan/CSO2-Server/web/register"
	"github.com/garyburd/redigo/redis"
	_ "github.com/mattn/go-sqlite3"
)

var (
	//SERVERVERSION 版本号
	SERVERVERSION = "v0.3.15"
	Redis         redis.Conn
)

func ReadHead(client net.Conn) ([]byte, bool) {
	SeqBuf := make([]byte, 1)
	headlen := HeaderLen - 1
	head, curlen := make([]byte, headlen), 0
	for {
		//DebugInfo(2, "Prepare", client.RemoteAddr().String())
		n, err := client.Read(SeqBuf)
		//DebugInfo(2, "Got", SeqBuf)
		if err != nil {
			return head, false
		}
		if n >= 1 && SeqBuf[0] == PacketTypeSignature {
			break
		}
		DebugInfo(2, "Recived a illegal head sig from", client.RemoteAddr().String())
	}
	for {
		//DebugInfo(2, "Prepare head", client.RemoteAddr().String())
		n, err := client.Read(head[curlen:])
		//DebugInfo(2, "Got", head)

		if err != nil {
			return head, false
		}
		curlen += n
		if curlen >= headlen {
			break
		}
	}
	return head, true
}

func ReadData(client net.Conn, len uint16) ([]byte, bool) {
	//DebugInfo(2, "Prepare data len:", len)
	data, curlen := make([]byte, len), 0
	for {
		//DebugInfo(2, "Prepare data", client.RemoteAddr().String())
		n, err := client.Read(data[curlen:])
		//DebugInfo(2, "Got", data)
		if err != nil {
			return data, false
		}
		curlen += n
		if curlen >= int(len) {
			break
		}
	}
	return data, true
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("检测到异常")
			fmt.Println("error:", err)
			fmt.Println("异常结束")
		}
	}()

	fmt.Println("Counter-Strike Online 2 Server", SERVERVERSION)
	fmt.Println("Initializing process ...")

	for k, v := range os.Args {
		if v == "-console" {
			os.Args = append(os.Args[:k], os.Args[k+1:]...)
			// 定义几个变量，用于接收命令行的参数值
			var user string
			var password string
			var host string
			var port string
			// &user 就是接收命令行中输入 -u 后面的参数值，其他同理
			flag.StringVar(&user, "username", "admin", "账号，默认为admin")
			flag.StringVar(&password, "password", "cso2server123", "密码，默认为cso2server123")
			flag.StringVar(&host, "ip", "localhost", "主机名，默认为localhost")
			flag.StringVar(&port, "port", "1315", "端口号，默认为1315")
			// 解析命令行参数写入注册的flag里
			flag.Parse()

			ToConsoleHost(user, password, host, port)

			return
		}
	}
	//get server exe path
	path, err := GetExePath()
	if err != nil {
		panic(err)
	}
	DBPath = path + "\\CSO2-Server\\database\\json\\"
	ReportPath = path + "\\CSO2-Server\\database\\report\\"

	//check folder
	checkFolder(path)

	//read configure
	Conf.InitConf(path)

	InitCSV(path)
	FullInventoryItem = CreateFullInventoryItem()
	FullInventoryReply = BuildFullInventoryInfo()
	DefaultInventoryItem = InitDefaultInventoryItem()
	InitBoxReply()
	InitCampaignReward()

	//read locales
	Locales.InitMotd(path)
	if Locales.InitLocales(path) {
		SetLocales()
	}

	//set val
	Level = Conf.DebugLevel
	LogFile = Conf.LogFile
	if Conf.MaxUsers <= 0 {
		MaxUsers = math.MaxUint32
	} else {
		MaxUsers = Conf.MaxUsers
	}

	//init Logger
	if LogFile != 0 {
		InitLoger(path, "CSO2-Server.log")
	}

	//init TCP
	server, err := net.Listen("tcp", fmt.Sprintf(":%d", Conf.PORT))
	if err != nil {
		fmt.Println("Init tcp socket error !\n")
		panic(err)
	}
	defer server.Close()

	//init UDP
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", Conf.HolePunchPort))
	if err != nil {
		fmt.Println("Init udp addr error !\n")
		panic(err)
	}
	holepunchserver, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Init udp socket error !\n")
		panic(err)
	}
	defer holepunchserver.Close()

	//Init Database
	if Conf.EnableDataBase != 0 {
		DB, err = InitDatabase(path + "\\CSO2-Server\\database\\sqlite\\cso2.db")
		if err != nil {
			fmt.Println("Init database failed !")
			Conf.EnableDataBase = 0
		} else {
			fmt.Println("Database connected !")
			defer DB.Close()
		}
	} else {
		DB = nil
	}

	//Init Redis
	if Conf.EnableRedis != 0 {
		Redis, err := InitRedis(Conf.RedisIP + ":" + strconv.Itoa(int(Conf.RedisPort)))
		if err != nil {
			fmt.Println("connect to redis server failed !")
			Conf.EnableRedis = 0
		} else {
			fmt.Println("Redis server connected !")
			defer Redis.Close()
		}
	}

	//Init converter
	InitConverter(Conf.CodePage)

	//Init MainServer Info
	MainServer = NewMainServer()
	InitExpTotal()

	//Start UDP Server
	go StartHolePunchServer(strconv.Itoa(int(Conf.HolePunchPort)), holepunchserver)

	//Start TCP Server
	go TCPServer(server)

	//Start BroadCast Service
	go BroadcastRoomList()

	//Start Register Server
	if Conf.EnableRegister != 0 {
		go OnRegister()
	}

	//Start console server
	if Conf.EnableConsole != 0 {
		go InitGMconsole()
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	_ = <-ch

	if SaveAllUsers() {
		fmt.Println("Save all users data success !")
	} else {
		fmt.Println("Save all users data failed !")
	}
	fmt.Println("Press CTRL+C again to close server")

	signal.Notify(ch, syscall.SIGINT)
	_ = <-ch
}

func TCPServer(server net.Listener) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("TCP server suffered a fault !")
			fmt.Println("error:", err)
			fmt.Println("Fault end!")
		}
	}()

	fmt.Println("Server is running at", "[AnyAdapter]:"+strconv.Itoa(int(Conf.PORT)))
	for {
		client, err := server.Accept()
		if err != nil {
			DebugInfo(2, "Server accept data error !\n")
			continue
		}
		DebugInfo(2, "Server accept a new connection request at", client.RemoteAddr().String())
		go RecvMessage(client)
	}
}

//RecvMessage 循环处理收到的包
func RecvMessage(client net.Conn) {
	var seq uint8 = 0
	var dataPacket PacketData

	defer client.Close() //关闭con
	defer func() {
		if err := recover(); err != nil {
			OnSendMessage(&seq, client, MessageDialogBox, GAME_SERVER_ERROR)
			fmt.Println("Client", client.RemoteAddr().String(), "suffered a fault !")
			fmt.Println(err)
			fmt.Println("Dump data", dataPacket.Data, "offset:", dataPacket.CurOffset)
			fmt.Println("Fault end!")
			OnLeaveRoom(client, true)
			DelUserWithConn(client)
		}
	}()

	client.Write([]byte("~SERVERCONNECTED\n"))

	for {
		//读取4字节数据包头部
		headBytes, err := ReadHead(client)
		if !err {
			goto close
		}
		var headPacket PacketHeader
		headPacket.Data = headBytes
		headPacket.PraseHeadPacket()
		// if !headPacket.IsGoodPacket {
		// 	DebugInfo(2, "Recived a illegal head from", client.RemoteAddr().String())
		// 	continue
		// }

		//读取数据部分
		bytes, err := ReadData(client, headPacket.Length)
		if !err {
			goto close
		}
		dataPacket = PacketData{
			bytes,
			headPacket.Sequence,
			headPacket.Length,
			bytes[0],
			1,
		}

		//执行功能
		switch dataPacket.Id {
		case PacketTypeQuickJoin:
			OnQuick(&dataPacket, client)
		case PacketTypeVersion:
			OnVersionPacket(&seq, client)
		case PacketTypeLogin:
			OnLogin(&seq, &dataPacket, client)
		case PacketTypeRequestChannels:
			OnServerList(client)
		case PacketTypeRequestRoomList:
			OnRoomList(&dataPacket, client)
		case PacketTypeRoom:
			OnRoomRequest(&dataPacket, client)
		case PacketTypeHost:
			OnHost(&dataPacket, client)
		case PacketTypeFavorite:
			OnFavorite(&dataPacket, client)
		case PacketTypeOption:
			OnOption(&dataPacket, client)
		case PacketTypePlayerInfo:
			OnPlayerInfo(&dataPacket, client)
		case PacketTypeChat:
			OnChat(&dataPacket, client)
		case PacketTypeAchievement:
			OnAchievement(&dataPacket, client)
		case PacketTypeAutomatch:
			OnAutoMatch(&dataPacket, client)
		case PacketTypeShop:
			OnShopRequest(&dataPacket, client)
		case PacketTypeReport:
			OnReportRequest(&dataPacket, client)
		case PacketTypeMail:
			OnMail(&dataPacket, client)
		case PacketTypeSupply:
			OnSupplyRequest(&dataPacket, client)
		case PacketTypePointLotto:
			OnPointLotto(&dataPacket, client)
		default:
			DebugInfo(2, "Unknown packet", dataPacket.Id, "from", client.RemoteAddr().String())
			//DebugInfo(2, "Unknown packet", dataPacket.Id, "from", client.RemoteAddr().String(), dataPacket.Data)
		}
	}

close:
	DebugInfo(1, "client", client.RemoteAddr().String(), "closed the connection")
	OnLeaveRoom(client, true)
	DelUserWithConn(client)
	return
}

func BroadcastRoomList() {
	for {
		timer := time.NewTimer(6 * time.Second)
		<-timer.C

		for _, v := range UsersManager.Users {
			if v != nil && v.CurrentChannelIndex > 0 && v.CurrentRoomId <= 0 {
				OnBroadcastRoomList(v.CurrentChannelServerIndex, v.CurrentChannelIndex, v)
			}
		}

	}
}

func checkFolder(path string) {
	rst, _ := PathExists(DBPath)
	if !rst {
		err := os.Mkdir(DBPath, os.ModePerm)
		if err != nil {
			fmt.Println("mkdir1 failed!", err)
		} else {
			fmt.Println("mkdir1 success!")
		}
	}

	folderpath := path + "\\CSO2-Server\\database\\report\\"
	rst, _ = PathExists(folderpath)
	if !rst {
		err := os.Mkdir(folderpath, os.ModePerm)
		if err != nil {
			fmt.Println("mkdir2 failed!", err)
		} else {
			fmt.Println("mkdir2 success!")
		}
	}
}

// func generate(path string) {
// 	file := path + "\\supplyList.csv"
// 	f, _ := os.Create(file)
// 	defer f.Close()
// 	f.WriteString(fmt.Sprintf("boxid,boxname,itemid,itemname,value\n"))
// 	for _, v := range BoxList {
// 		for _, item := range v.Items {
// 			if len(ItemList[item.ItemID].Name) <= 0 {
// 				continue
// 			}
// 			f.WriteString(fmt.Sprintf("%d,%s,%d,%s,%d\n", v.BoxID, v.BoxName, item.ItemID, ItemList[item.ItemID].Name, item.Value))
// 		}
// 	}
// }
