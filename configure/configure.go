package configure

import (
	"fmt"
	"io/ioutil"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type CSO2Conf struct {
	PORT             uint32
	HolePunchPort    uint32
	EnableRedis      uint32
	EnableDataBase   uint32
	MaxUsers         uint32
	EnableShop       uint32
	UnlockAllWeapons uint32
	UnlockAllSkills  uint32
	RedisIP          string
	RedisPort        uint32
	DebugLevel       uint32
	LogFile          uint32
	EnableRegister   uint32
	EnableMail       uint32
	REGPort          uint32
	REGEmail         string
	REGPassWord      string
	REGSMTPaddr      string
	LocaleFile       string
	CodePage         string
	EnableConsole    uint32
	GMport           uint32
	GMusername       string
	GMpassword       string
}

type CSO2Locales struct {
	GAME_ROOM_LEAVE_EARLY       string
	GAME_SERVER_ERROR           string
	GAME_LOGIN_ALREADY          string
	GAME_LOGIN_ERROR            string
	GAME_ROOM_COUNT_MODE_ERROR  string
	GAME_ROOM_JOIN_ERROR        string
	GAME_GM_ADD_ALLWEAPONS      string
	GAME_CHANNEL_MESSAGE_NOT_IN string
	GAME_GM_NO_AUTHORIZE        string
	GAME_USER_NEW_ITEM          string
	GAME_USER_NEW_ITEM_RESTART  string
	GAME_CHANNEL_MESSAGE        string

	MAIL_EMPTY       string
	MAIL_ERROR       string
	MAIL_SENT        string
	USERNAME_EMPTY   string
	USERNAME_EXISTS  string
	GAMENAME_EMPTY   string
	GAMENAME_EXISTS  string
	PASSWORD_EMPTY   string
	CODE_EMPTY       string
	CODE_WRONG       string
	NAME_ERROR       string
	DATABASE_ERROR   string
	REGISTER_SUCCESS string

	MOTD []byte
}

var (
	Conf    CSO2Conf
	Locales CSO2Locales

	MAIL_EMPTY       = "提示：邮箱不能为空！"
	MAIL_ERROR       = "提示：请输入正确的邮箱！"
	MAIL_SENT        = "已发送，请在一分钟之内完成注册！"
	USERNAME_EMPTY   = "提示：用户名不能为空！"
	USERNAME_EXISTS  = "提示：用户名已存在！"
	GAMENAME_EMPTY   = "提示：游戏昵称不能为空！"
	GAMENAME_EXISTS  = "提示：游戏昵称已存在！"
	PASSWORD_EMPTY   = "提示：密码不能为空！"
	CODE_EMPTY       = "提示：验证码不能为空！"
	CODE_WRONG       = "提示：验证码不正确！"
	NAME_ERROR       = "提示：用户名或昵称含有非法字符！"
	DATABASE_ERROR   = "提示：数据库错误,注册失败！"
	REGISTER_SUCCESS = "注册成功!"
)

func (conf *CSO2Conf) InitConf(path string) {
	if conf == nil {
		return
	}
	fmt.Printf("Reading configure file ...\n")
	ini_parser := IniParser{}
	file := path + "\\CSO2-Server\\configure\\server.conf"
	if err := ini_parser.LoadIni(file); err != nil {
		fmt.Printf("Loading config file error[%s]\n", err.Error())
		fmt.Printf("Using default data ...\n")
		conf.EnableRedis = 0
		conf.EnableDataBase = 1
		conf.MaxUsers = 0
		conf.EnableShop = 0
		conf.UnlockAllWeapons = 1
		conf.UnlockAllSkills = 1
		conf.PORT = 30001
		conf.HolePunchPort = 30002
		conf.RedisIP = "127.0.0.1"
		conf.RedisPort = 6379
		conf.DebugLevel = 2
		conf.LogFile = 1
		conf.EnableRegister = 1
		conf.EnableMail = 0
		Conf.LocaleFile = "zh-cn.ini"
		Conf.CodePage = "gbk"
		conf.EnableConsole = 0
		return
	}
	conf.EnableRedis = ini_parser.IniGetUint32("Database", "EnableRedis")
	conf.EnableDataBase = ini_parser.IniGetUint32("Database", "EnableDataBase")
	conf.MaxUsers = ini_parser.IniGetUint32("Server", "MaxUsers")
	if conf.MaxUsers < 0 {
		conf.MaxUsers = 0
	}
	conf.EnableShop = ini_parser.IniGetUint32("Server", "EnableShop")
	conf.UnlockAllWeapons = ini_parser.IniGetUint32("Server", "UnlockAllWeapons")
	conf.UnlockAllSkills = ini_parser.IniGetUint32("Server", "UnlockAllSkills")
	conf.PORT = ini_parser.IniGetUint32("Server", "TCPPort")
	conf.HolePunchPort = ini_parser.IniGetUint32("Server", "UDPPort")
	conf.RedisIP = ini_parser.IniGetString("Server", "RedisIP")
	conf.RedisPort = ini_parser.IniGetUint32("Server", "RedisPort")
	conf.DebugLevel = ini_parser.IniGetUint32("Debug", "DebugLevel")
	if conf.DebugLevel > 2 || conf.DebugLevel < 0 {
		conf.DebugLevel = 2
	}
	conf.LogFile = ini_parser.IniGetUint32("Debug", "LogFile")
	conf.EnableRegister = ini_parser.IniGetUint32("Register", "EnableRegister")
	conf.EnableMail = ini_parser.IniGetUint32("Register", "EnableMail")
	conf.REGPort = ini_parser.IniGetUint32("Register", "REGPort")
	conf.REGEmail = ini_parser.IniGetString("Register", "REGEmail")
	conf.REGPassWord = ini_parser.IniGetString("Register", "REGPassWord")
	conf.REGSMTPaddr = ini_parser.IniGetString("Register", "REGSMTPaddr")
	Conf.LocaleFile = ini_parser.IniGetString("Locale", "LocaleFile")
	Conf.CodePage = ini_parser.IniGetString("Encode", "CodePage")
	conf.EnableConsole = ini_parser.IniGetUint32("Console", "EnableConsole")
	conf.GMport = ini_parser.IniGetUint32("Console", "GMport")
	Conf.GMusername = ini_parser.IniGetString("Console", "GMusername")
	Conf.GMpassword = ini_parser.IniGetString("Console", "GMpassword")
}

func (locales *CSO2Locales) InitLocales(path string) bool {
	if locales == nil {
		return false
	}
	fmt.Printf("Reading locale < " + Conf.LocaleFile + " > ...\n")
	ini_parser := IniParser{}
	file := path + "\\CSO2-Server\\locales\\" + Conf.LocaleFile
	if err := ini_parser.LoadIni(file); err != nil {
		fmt.Printf("Loading locale file error[%s]\n", err.Error())
		fmt.Printf("Using default data ...\n")
		return false
	}
	locales.GAME_ROOM_LEAVE_EARLY = ini_parser.IniGetString("System", "GAME_ROOM_LEAVE_EARLY")
	locales.GAME_SERVER_ERROR = ini_parser.IniGetString("System", "GAME_SERVER_ERROR")
	locales.GAME_LOGIN_ALREADY = ini_parser.IniGetString("System", "GAME_LOGIN_ALREADY")
	locales.GAME_LOGIN_ERROR = ini_parser.IniGetString("System", "GAME_LOGIN_ERROR")
	locales.GAME_ROOM_COUNT_MODE_ERROR = ini_parser.IniGetString("System", "GAME_ROOM_COUNT_MODE_ERROR")
	locales.GAME_ROOM_JOIN_ERROR = ini_parser.IniGetString("System", "GAME_ROOM_JOIN_ERROR")
	locales.GAME_GM_ADD_ALLWEAPONS = ini_parser.IniGetString("System", "GAME_GM_ADD_ALLWEAPONS")
	locales.GAME_CHANNEL_MESSAGE_NOT_IN = ini_parser.IniGetString("System", "GAME_CHANNEL_MESSAGE_NOT_IN")
	locales.GAME_GM_NO_AUTHORIZE = ini_parser.IniGetString("System", "GAME_GM_NO_AUTHORIZE")
	locales.GAME_USER_NEW_ITEM = ini_parser.IniGetString("System", "GAME_USER_NEW_ITEM")
	locales.GAME_USER_NEW_ITEM_RESTART = ini_parser.IniGetString("System", "GAME_USER_NEW_ITEM_RESTART")
	locales.GAME_CHANNEL_MESSAGE = ini_parser.IniGetString("System", "GAME_CHANNEL_MESSAGE")

	locales.MAIL_EMPTY = ini_parser.IniGetString("Register", "MAIL_EMPTY")
	locales.MAIL_ERROR = ini_parser.IniGetString("Register", "MAIL_ERROR")
	locales.USERNAME_EMPTY = ini_parser.IniGetString("Register", "USERNAME_EMPTY")
	locales.USERNAME_EXISTS = ini_parser.IniGetString("Register", "USERNAME_EXISTS")
	locales.GAMENAME_EMPTY = ini_parser.IniGetString("Register", "GAMENAME_EMPTY")
	locales.GAMENAME_EXISTS = ini_parser.IniGetString("Register", "GAMENAME_EXISTS")
	locales.PASSWORD_EMPTY = ini_parser.IniGetString("Register", "PASSWORD_EMPTY")
	locales.CODE_EMPTY = ini_parser.IniGetString("Register", "CODE_EMPTY")
	locales.CODE_WRONG = ini_parser.IniGetString("Register", "CODE_WRONG")
	locales.NAME_ERROR = ini_parser.IniGetString("Register", "NAME_ERROR")
	locales.DATABASE_ERROR = ini_parser.IniGetString("Register", "DATABASE_ERROR")
	locales.REGISTER_SUCCESS = ini_parser.IniGetString("Register", "REGISTER_SUCCESS")
	return true
}

func (locales *CSO2Locales) InitMotd(path string) {
	if locales == nil {
		return
	}
	fmt.Printf("Reading motd ...\n")
	filepath := path + "\\CSO2-Server\\locales\\motd.txt"
	dataEncoded, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Using default motd ...\n")
		locales.MOTD = []byte("You are playing Counter-Strike Online 2 on CSO2-Server.Please delete in 24 hours after loading !Visit the github web site github.com/KouKouChan/CSO2-Server")
		return
	}
	locales.MOTD = dataEncoded
	return
}

func SetLocales() {
	fmt.Printf("Setting locale < " + Conf.LocaleFile + " > ...\n")
	GAME_ROOM_LEAVE_EARLY = []byte(Locales.GAME_ROOM_LEAVE_EARLY)
	GAME_SERVER_ERROR = []byte(Locales.GAME_SERVER_ERROR)
	GAME_LOGIN_ALREADY = []byte(Locales.GAME_LOGIN_ALREADY)
	GAME_LOGIN_ERROR = []byte(Locales.GAME_LOGIN_ERROR)
	GAME_ROOM_COUNT_MODE_ERROR = []byte(Locales.GAME_ROOM_COUNT_MODE_ERROR)
	GAME_ROOM_JOIN_ERROR = []byte(Locales.GAME_ROOM_JOIN_ERROR)
	GAME_GM_ADD_ALLWEAPONS = []byte(Locales.GAME_GM_ADD_ALLWEAPONS)
	GAME_CHANNEL_MESSAGE_NOT_IN = []byte(Locales.GAME_CHANNEL_MESSAGE_NOT_IN)
	GAME_GM_NO_AUTHORIZE = []byte(Locales.GAME_GM_NO_AUTHORIZE)
	GAME_USER_NEW_ITEM = []byte(Locales.GAME_USER_NEW_ITEM)
	GAME_USER_NEW_ITEM_RESTART = []byte(Locales.GAME_USER_NEW_ITEM_RESTART)
	GAME_CHANNEL_MESSAGE = Locales.GAME_CHANNEL_MESSAGE

	MAIL_EMPTY = Locales.MAIL_EMPTY
	MAIL_ERROR = Locales.MAIL_ERROR
	MAIL_SENT = Locales.MAIL_SENT
	USERNAME_EMPTY = Locales.USERNAME_EMPTY
	USERNAME_EXISTS = Locales.USERNAME_EXISTS
	GAMENAME_EMPTY = Locales.GAMENAME_EMPTY
	GAMENAME_EXISTS = Locales.GAMENAME_EXISTS
	PASSWORD_EMPTY = Locales.PASSWORD_EMPTY
	CODE_EMPTY = Locales.CODE_EMPTY
	CODE_WRONG = Locales.CODE_WRONG
	NAME_ERROR = Locales.NAME_ERROR
	DATABASE_ERROR = Locales.DATABASE_ERROR
	REGISTER_SUCCESS = Locales.REGISTER_SUCCESS
}
