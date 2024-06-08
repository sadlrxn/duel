package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/models"
)

type ipGeolocationResponse struct {
	City          string `json:"city"`
	RegionISOCode string `json:"region_iso_code"`
	CountryCode   string `json:"country_code"`
}

func CheckIpAddress(blackList []string, address string) bool {
	// isValidIp, _ := regexp.MatchString(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`, address)
	// if !isValidIp {
	// 	return false
	// }
	// for i := 0; i < len(blackList); i++ {
	// 	blackIp := blackList[i]
	// 	segs := strings.Split(blackIp, ".")
	// 	pattern := "^"
	// 	for j := 0; j < len(segs); j++ {
	// 		if _, err := strconv.Atoi(segs[j]); err == nil {
	// 			pattern += segs[j]
	// 		} else {
	// 			pattern += `\d+`
	// 		}
	// 		if j < len(segs)-1 {
	// 			pattern += `\.`
	// 		}
	// 	}
	// 	pattern += "$"

	// 	isBlacklistedIp, _ := regexp.MatchString(pattern, address)
	// 	if isBlacklistedIp {
	// 		return false
	// 	}
	// }
	return true
}

func CheckIpRegionWithAbstractAPI(address string, regions map[string]bool) (string, bool) {
	config := config.Get()
	url := fmt.Sprintf(`https://ipgeolocation.abstractapi.com/v1/?api_key=%s&fields=country_code,region_iso_code,city&ip_address=%s`, config.IpGeolocationApiKey, address)
	res, err := http.Get(url)
	if err != nil {
		return "", false
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", false
	}

	var resObj ipGeolocationResponse
	json.Unmarshal(body, &resObj)

	return resObj.CountryCode, !regions[resObj.CountryCode]
}

func CheckBannedUser(userID uint) bool {
	db := db.GetDB()
	var user models.User
	db.First(&user, userID)
	return user.Banned
}
