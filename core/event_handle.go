package core

import (
	"encoding/json"
	"log"
)

//PlayerInfo ...
type PlayerInfo struct {
	UID          int
	Username     string
	Imm          int
	Spt          bool
	Vip          bool
	Ctb          bool
	Opt          bool
	Adm          bool
	Own          bool
	Tviplevel    int
	Grp          int
	OnlineTotal  int
	OnlineToday  int
	OnlineOB     int
	OnlinePlay   int
	ConnectTimes int
	Vitality     int
	TrackingID   int
}

//PlayerConnection ...
type PlayerConnection struct {
	SteamID     string `json:"SteamID"`
	CIndex      int    `json:"CIndex"`
	IP          string `json:"IP"`
	JoinTime    int    `json:"JoinTime"`
	TodayDate   int    `json:"TodayDate"`
	Map         string `json:"Map"`
	ServerID    int    `json:"ServerID"`
	ServerModID int    `json:"ServerModID"`
}

//Chat ...
type Chat struct {
	ServerID    int    `json:"ServerID"`
	ServerModID int    `json:"ServerModID"`
	PlayerName  string `json:"PlayerName"`
	Msg         string `json:"Msg"`
}

//KillStats ...
type KillStats struct {
	Attacker int    `json:"Attacker"`
	Assister int    `json:"Assister"`
	Victim   int    `json:"Victim"`
	Weapon   string `json:"Weapon"`
	Headshot bool   `json:"Headshot"`
}

//PlayerDisconnection ...
type PlayerDisconnection struct {
	UID int `json:"UID"`
}

//UserStats ...
type UserStats struct {
	UID         int    `json:"UID"`
	SessionID   int    `json:"SessionID"`
	TodayOnline int    `json:"TodayOnline"`
	TotalOnline int    `json:"TotalOnline"`
	SpecOnline  int    `json:"SpecOnline"`
	PlayOnline  int    `json:"PlayOnline"`
	UserName    string `json:"UserName"`
}

//EventData ...
type EventData struct {
	Event               string              `json:"Event"`
	PlayerConnection    PlayerConnection    `json:"PlayerConnection"`
	AllServersChat      Chat                `json:"AllServersChat"`
	KillStats           KillStats           `json:"KillStats"`
	PlayerDisconnection PlayerDisconnection `json:"PlayerDisconnection"`
	UserStats           UserStats           `json:"UserStats"`
	SQLSave             string              `json:"SQLSave"`
}

//EventHandle ...
func EventHandle(msg string, serNum int) {
	data := EventData{}
	err := json.Unmarshal([]byte(msg), &data)
	if err != nil {
		log.Println("Json解析错误: ", err, msg)
		return
	}

	switch {
	case data.Event == "AllServersChat":
		AllChatHandle(data, serNum)

	case data.Event == "PlayerConnection":
		PlayerConnHandle(data, serNum)

	case data.Event == "SQLSave":
		SQLSaveHandle(data, serNum)

	//case data.Event == "KillStats":
	//SQLSaveHandle(data, serNum)

	//case data.Event == "PlayerDisconnection":
	//SQLSaveHandle(data, serNum)

	case data.Event == "UserStats":
		StatsHandle(data, serNum)

	case data.Event == "RELOADSETTING":
		ReloadSetting()
	}
}

//SQLSaveHandle ...
func SQLSaveHandle(data EventData, serNum int) {
	_, err := db.Exec(data.SQLSave)
	if !CheckError(err) {
		log.Println("数据", data)
	}
}

//PlayerConnHandle ...
func PlayerConnHandle(data EventData, serNum int) {
	playerinfo := data.PlayerConnection
	var player PlayerInfo

	row, err := JoinQuery.Query(playerinfo.SteamID, playerinfo.ServerID, playerinfo.ServerModID, playerinfo.IP, playerinfo.Map, playerinfo.JoinTime, playerinfo.TodayDate)

	if !CheckError(err) {
		log.Println("数据", data)
		return
	}

	row.Next()
	row.Scan(&player.UID, &player.Username, &player.Imm, &player.Spt, &player.Vip, &player.Ctb, &player.Opt, &player.Adm, &player.Own, &player.Tviplevel, &player.Grp, &player.OnlineTotal, &player.OnlineToday, &player.OnlineOB, &player.OnlinePlay, &player.ConnectTimes, &player.Vitality, &player.TrackingID)
	row.Close()

	buff := struct {
		Event      string     `json:"Event"`
		PlayerInfo PlayerInfo `json:"PlayerInfo"`
		CIndex     int        `json:"CIndex"`
		SteamID    string     `json:"SteamID"`
	}{"PlayerInfo", player, playerinfo.CIndex, playerinfo.SteamID}

	json, _ := json.Marshal(buff)

	sersChan[serNum] <- string(json)

	return
}

//AllChatHandle ...
func AllChatHandle(data EventData, serNum int) {
	buff := struct {
		Event          string `json:"Event"`
		AllServersChat Chat   `json:"AllServersChat"`
	}{"AllServersChat", data.AllServersChat}

	json, _ := json.Marshal(buff)

	for k := range sersChan {
		sersChan[k] <- string(json)
	}
}

//StatsHandle ...
func StatsHandle(data EventData, serNum int) {
	_, err := StatsQuery.Exec(data.UserStats.UID, data.UserStats.SessionID, data.UserStats.TodayOnline, data.UserStats.TotalOnline, data.UserStats.SpecOnline, data.UserStats.PlayOnline, data.UserStats.UserName)
	if !CheckError(err) {
		log.Println("数据", data)
	}
}
