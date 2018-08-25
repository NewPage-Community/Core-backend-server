package core

import (
	"encoding/json"
	"time"
)

// CheckBan ...
func CheckBan(steamid string, serverid int, modid int, ip string) (bool, int, int64, string, string) {
	row, err := db.Query("SELECT `bid`, `bantype`, `sid`, `mid`, `ends`, `adminname`, `reason` FROM `np_bans` WHERE `steamid` = ? AND `bRemovedBy` = -1 AND (`ends` > ? OR `length` = 0) ORDER BY `created` DESC", steamid, time.Now().Unix())

	if !CheckError(err) {
		return false, 0, 0, "", ""
	}

	for row.Next() {
		var banType, sid, mid, bid int
		var reason, adminname string
		var ends int64

		row.Scan(&bid, &banType, &sid, &mid, &ends, &adminname, &reason)

		if (banType == 1 && mid == modid) || (banType == 2 && sid == serverid) || banType == 0 {
			db.Exec("INSERT INTO `np_blocks` VALUES (DEFAULT, ?, ?, ?)", bid, ip, time.Now().Unix())
			row.Close()
			return true, banType, ends, reason, adminname
		}
	}

	row.Close()
	return false, 0, 0, "", ""
}

// AddBan ...
func AddBan(data EventData, serNum int) {
	if data.BanInfo.UID == -1 {
		row, err := db.Query("SELECT `uid`, `username` FROM `np_users` WHERE `steamid` = '?'", data.BanInfo.SteamID)

		if CheckError(err) {
			row.Next()
			row.Scan(&data.BanInfo.UID, &data.BanInfo.NikeName)
			row.Close()
		}
	}

	var ETime int64

	if data.BanInfo.Length == 0 {
		ETime = 0
	} else {
		ETime = time.Now().Unix() + int64(data.BanInfo.Length)
	}

	_, err := db.Exec("INSERT INTO `np_bans` VALUES (DEFAULT, ?, '?', '?', '?', ?, ?, ?, ?, ?, ?, ?, '?', '?', '', -1)", data.BanInfo.UID, data.BanInfo.SteamID, data.BanInfo.IP, data.BanInfo.NikeName, time.Now().Unix(), data.BanInfo.Length, ETime, data.BanInfo.BanType, data.BanInfo.ServerID, data.BanInfo.ServerModID, data.BanInfo.AdminID, data.BanInfo.AdminName, data.BanInfo.Reason)

	if !CheckError(err) {
		return
	}

	banInfo := BanClient{data.BanInfo.SteamID, data.BanInfo.Length, time.Now().Unix() + int64(data.BanInfo.Length), data.BanInfo.BanType, data.BanInfo.AdminID, data.BanInfo.AdminName, data.BanInfo.Reason}

	buff := struct {
		Event     string    `json:"Event"`
		BanClient BanClient `json:"BanClient"`
	}{"BanClient", banInfo}

	json, _ := json.Marshal(buff)
	sersChan[serNum] <- string(json)
}
