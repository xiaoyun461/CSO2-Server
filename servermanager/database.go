package servermanager

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

var (
	DB     *sql.DB
	Dblock sync.Mutex
	DBPath string
)

//从数据库中读取用户数据
//如果是新用户则保存到数据库中
func GetUserFromDatabase(loginname string, passwd []byte) (*User, int) {
	if DB != nil {
		//query, err := DB.Prepare("SELECT * FROM userinfo WHERE LoginName = ?")
		//if err == nil {
		filepath := DBPath + loginname
		rst, err := PathExists(filepath)
		if rst {
			//defer query.Close()
			u := GetNewUser()
			//var inventory []byte
			//var clanID uint32
			Dblock.Lock()
			dataEncoded, _ := ioutil.ReadFile(filepath)
			// err = query.QueryRow(loginname).Scan(&u.UserName, &u.IngameName, &u.Password, &u.Level, &u.Rank,
			// 	&u.RankFrame, &u.Points, &u.CurrentExp, &u.PlayedMatches, &u.Wins, &u.Kills,
			// 	&u.Headshots, &u.Deaths, &u.Assists, &u.Accuracy, &u.SecondsPlayed, &u.NetCafeName,
			// 	&u.Cash, &clanID, &u.WorldRank, &u.Mpoints, &u.TitleId, &u.UnlockedTitles, &u.Signature,
			// 	&u.BestGamemode, &u.BestMap, &u.UnlockedAchievements, &u.Avatar, &u.UnlockedAvatars,
			// 	&u.VipLevel, &u.VipXp, &u.SkillHumanCurXp, &u.SkillHumanPoints, &u.SkillZombieCurXp,
			// 	&u.SkillZombiePoints, &inventory, &u.UserMail)

			Dblock.Unlock()
			err = json.Unmarshal(dataEncoded, &u)
			if err != nil {
				DebugInfo(1, "Suffered a error while getting User", string(loginname)+"'s data !", err)
				return nil, USER_UNKOWN_ERROR
			}
			//检查密码
			str := fmt.Sprintf("%x", md5.Sum([]byte(string(loginname)+string(passwd))))
			for i := 0; i < 16; i++ {
				if str[i] != u.Password[i] {
					//u = GetNewUser()
					//u.SetID(0)
					return nil, USER_PASSWD_ERROR
				}
			}
			// //设置仓库
			// u.Inventory = praseInventory(inventory)
			//设置战队...
			DebugInfo(1, "User", string(u.UserName), "data found !")
			u.SetID(GetNewUserID())
			//u.MaxExp = LevelExp[u.Level-1]
			return &u, USER_LOGIN_SUCCESS
		} else {
			return nil, USER_NOT_FOUND
		}
	}
	u := GetNewUser()
	u.SetID(GetNewUserID())
	u.SetUserName([]byte(loginname), []byte(loginname))
	u.Password = passwd
	return &u, USER_LOGIN_SUCCESS
}

func praseInventory(inventory []byte) UserInventory {
	var inv UserInventory
	offset := 0
	inv.NumOfItem = ReadUint16(inventory, &offset)
	for i := 0; i < int(inv.NumOfItem); i++ {
		var it UserInventoryItem
		it.Id = ReadUint32(inventory, &offset)
		it.Count = ReadUint16(inventory, &offset)
		inv.Items = append(inv.Items, it)
	}
	inv.CTModel = ReadUint32(inventory, &offset)
	inv.TModel = ReadUint32(inventory, &offset)
	inv.HeadItem = ReadUint32(inventory, &offset)
	inv.GloveItem = ReadUint32(inventory, &offset)
	inv.BackItem = ReadUint32(inventory, &offset)
	inv.StepsItem = ReadUint32(inventory, &offset)
	inv.CardItem = ReadUint32(inventory, &offset)
	inv.SprayItem = ReadUint32(inventory, &offset)
	//buymenu
	len := ReadUint8(inventory, &offset)
	inv.BuyMenu.Pistols = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.BuyMenu.Shotguns = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.BuyMenu.Smgs = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.BuyMenu.Rifles = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.BuyMenu.Snipers = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.BuyMenu.Machineguns = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.BuyMenu.Melees = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.BuyMenu.Equipment = ReadUint32Array(inventory, &offset, int(len))
	//loadouts
	len = ReadUint8(inventory, &offset)
	for i := 0; i < int(len); i++ {
		var ul UserLoadout
		l := ReadUint8(inventory, &offset)
		ul.Items = ReadUint32Array(inventory, &offset, int(l))
		inv.Loadouts = append(inv.Loadouts, ul)
	}
	return inv
}

func InventoryToBytes(inventory *UserInventory) []byte {
	buf := make([]byte, 8096)
	offset := 0
	WriteUint16(&buf, inventory.NumOfItem, &offset)
	for i := 0; i < int(inventory.NumOfItem); i++ {
		WriteUint32(&buf, inventory.Items[i].Id, &offset)
		WriteUint16(&buf, inventory.Items[i].Count, &offset)
	}
	WriteUint32(&buf, inventory.CTModel, &offset)
	WriteUint32(&buf, inventory.TModel, &offset)
	WriteUint32(&buf, inventory.HeadItem, &offset)
	WriteUint32(&buf, inventory.GloveItem, &offset)
	WriteUint32(&buf, inventory.BackItem, &offset)
	WriteUint32(&buf, inventory.StepsItem, &offset)
	WriteUint32(&buf, inventory.CardItem, &offset)
	WriteUint32(&buf, inventory.SprayItem, &offset)
	//buymenu
	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Pistols)), &offset)
	WriteUint32Array(&buf, inventory.BuyMenu.Pistols, &offset)
	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Shotguns)), &offset)
	WriteUint32Array(&buf, inventory.BuyMenu.Shotguns, &offset)
	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Smgs)), &offset)
	WriteUint32Array(&buf, inventory.BuyMenu.Smgs, &offset)
	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Rifles)), &offset)
	WriteUint32Array(&buf, inventory.BuyMenu.Rifles, &offset)
	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Snipers)), &offset)
	WriteUint32Array(&buf, inventory.BuyMenu.Snipers, &offset)
	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Machineguns)), &offset)
	WriteUint32Array(&buf, inventory.BuyMenu.Machineguns, &offset)
	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Melees)), &offset)
	WriteUint32Array(&buf, inventory.BuyMenu.Melees, &offset)
	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Equipment)), &offset)
	WriteUint32Array(&buf, inventory.BuyMenu.Equipment, &offset)
	//loadouts
	WriteUint8(&buf, uint8(len(inventory.Loadouts)), &offset)
	for i := 0; i < len(inventory.Loadouts); i++ {
		WriteUint8(&buf, uint8(len(inventory.Loadouts[i].Items)), &offset)
		WriteUint32Array(&buf, inventory.Loadouts[i].Items, &offset)
	}
	return buf[:offset]
}

func AddUserToDB(u *User) error {
	if DB == nil {
		return errors.New("DataBase not connected")
	}
	// stmt, err := DB.Prepare(`INSERT INTO userinfo(LoginName, UserName, PassWord,
	// 	Level, Rank, RankFrame, Points, CurrentExp, PlayedMatches, Wins, Kills,
	// 	HeadShots, Deathes, Assists, accuracy, SecondsPlayed, netCafeName,
	// 	Cash, ClanID, WorldRank, Mpoints, TitleID, UnlockefTitleID, signature,
	// 	bestGamemode, bestMap, unlockedAchievements, avatar, unlockedAvatars,
	// 	viplevel, vipXp, skillHumanCurXp, skillHumanPoints, skillZombieCurXp,
	// 	skillZombiePoints, Inventory, UserMail) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?
	// 	   ,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`) //36个
	// if err != nil {
	// 	return err
	// }
	// defer stmt.Close()

	filepath := DBPath + string(u.UserName)
	data, _ := json.MarshalIndent(u, "", "     ")
	Dblock.Lock()
	err := ioutil.WriteFile(filepath, data, 0644)

	// _, err = stmt.Exec(u.UserName, u.IngameName, u.Password, u.Level, u.Rank,
	// 	u.RankFrame, u.Points, u.CurrentExp, u.PlayedMatches, u.Wins, u.Kills,
	// 	u.Headshots, u.Deaths, u.Assists, u.Accuracy, u.SecondsPlayed, u.NetCafeName,
	// 	u.Cash, 0, u.WorldRank, u.Mpoints, u.TitleId, u.UnlockedTitles, u.Signature,
	// 	u.BestGamemode, u.BestMap, u.UnlockedAchievements, u.Avatar, u.UnlockedAvatars,
	// 	u.VipLevel, u.VipXp, u.SkillHumanCurXp, u.SkillHumanPoints, u.SkillZombieCurXp,
	// 	u.SkillZombiePoints, InventoryToBytes(&u.Inventory), u.UserMail)

	Dblock.Unlock()
	if err != nil {
		return err
	}

	filepath = DBPath + string(u.IngameName) + ".check"
	Dblock.Lock()
	err = ioutil.WriteFile(filepath, u.IngameName, 0644)
	Dblock.Unlock()
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserToDB(u *User) error {
	if DB == nil {
		return errors.New("DataBase not connected,can't save user's data !")
	}
	// stmt, err := DB.Prepare(`Update userinfo set Level=?,
	// 	Rank=?, RankFrame=?, Points=?, CurrentExp=?, PlayedMatches=?, Wins=?, Kills=?,
	// 	HeadShots=?, Deathes=?, Assists=?, accuracy=?, SecondsPlayed=?, netCafeName=?,
	// 	Cash=?, ClanID=?, WorldRank=?, Mpoints=?, TitleID=?, UnlockefTitleID=?, signature=?,
	// 	bestGamemode=?, bestMap=?, unlockedAchievements=?, avatar=?, unlockedAvatars=?,
	// 	viplevel=?, vipXp=?, skillHumanCurXp=?, skillHumanPoints=?, skillZombieCurXp=?,
	// 	skillZombiePoints=?, Inventory=? WHERE LoginName=? `) //36个
	// if err != nil {
	// 	return err
	// }
	// defer stmt.Close()

	filepath := DBPath + string(u.UserName)
	data, _ := json.MarshalIndent(u, "", "     ")
	Dblock.Lock()
	err := ioutil.WriteFile(filepath, data, 0644)

	// _, err = stmt.Exec(u.Level, u.Rank,
	// 	u.RankFrame, u.Points, u.CurrentExp, u.PlayedMatches, u.Wins, u.Kills,
	// 	u.Headshots, u.Deaths, u.Assists, u.Accuracy, u.SecondsPlayed, u.NetCafeName,
	// 	u.Cash, 0, u.WorldRank, u.Mpoints, u.TitleId, u.UnlockedTitles, u.Signature,
	// 	u.BestGamemode, u.BestMap, u.UnlockedAchievements, u.Avatar, u.UnlockedAvatars,
	// 	u.VipLevel, u.VipXp, u.SkillHumanCurXp, u.SkillHumanPoints, u.SkillZombieCurXp,
	// 	u.SkillZombiePoints, InventoryToBytes(&u.Inventory), u.UserName)

	Dblock.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func IsExistsMail(mail []byte) bool {
	if DB != nil {
		query, err := DB.Prepare("SELECT * FROM userinfo WHERE UserMail = ?")
		if err == nil {
			defer query.Close()
			Dblock.Lock()
			rows, err := query.Query(mail)
			Dblock.Unlock()
			if err != nil {
				return false
			}
			defer rows.Close()
			if rows.Next() {
				return true
			}
		}
		//存在风险，如果出错时候其实该用户存在，那么会出现冗余
		return false
	}
	return false
}

func IsExistsUser(username []byte) bool {
	if DB != nil {
		//query, err := DB.Prepare("SELECT * FROM userinfo WHERE LoginName = ?")
		filepath := DBPath + string(username)
		rst, _ := PathExists(filepath)
		if rst {
			return true
			// defer query.Close()
			// dblock.Lock()
			// rows, err := query.Query(username)
			// dblock.Unlock()
			// if err != nil {
			// 	return false
			// }
			// defer rows.Close()
			// if rows.Next() {
			// 	return true
			// }
		}
		return false
	}
	return false
}

func IsExistsIngameName(name []byte) bool {
	if DB != nil {
		//query, err := DB.Prepare("SELECT * FROM userinfo WHERE LoginName = ?")
		filepath := DBPath + string(name) + ".check"
		rst, _ := PathExists(filepath)
		if rst {
			return true
			// defer query.Close()
			// dblock.Lock()
			// rows, err := query.Query(username)
			// dblock.Unlock()
			// if err != nil {
			// 	return false
			// }
			// defer rows.Close()
			// if rows.Next() {
			// 	return true
			// }
		}
		return false
	}
	return false
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
