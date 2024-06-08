package user

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"regexp"
	"strings"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/utils"
)

/*
* @Ixternal
* Generate a random user name for a new user.
 */
func generateUserName() (string, error) {
	// 1. Send request to generate random user name.
	var url = "https://namey.muffinlabs.com/name.json?count=1&type=surname&frequency=common"
	res, err := http.Get(url)
	if err != nil {
		return "", utils.MakeError(
			"user_controller",
			"generateUserName",
			"failed to request random user name",
			err,
		)
	}

	// 2. Read response body as byte array.
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", utils.MakeError(
			"user_controller",
			"generateUserName",
			"failed to read response body",
			err,
		)
	}

	// 3. Parse response body to string array.
	var resObj []string
	if err := json.Unmarshal(body, &resObj); err != nil {
		return "", utils.MakeError(
			"user_controller",
			"generateUserName",
			"failed to parse response to json object",
			err,
		)
	}

	// 4. Validate response.
	if len(resObj) == 0 {
		return "", utils.MakeError(
			"user_controller",
			"generateUserName",
			"invalid response",
			fmt.Errorf("response: %v", resObj),
		)
	}

	// 5. Cut name if longer than 12 characters.
	name := resObj[0]
	if len(name) > 12 {
		name = name[:12]
	}

	// 6. Count users with similar names.
	count := countSimilarNames(name)
	if count == -1 {
		return "", utils.MakeError(
			"user_controller",
			"generateUserName",
			"failed to count users with similar names",
			fmt.Errorf(
				"name: %s", name,
			),
		)
	}

	// 7. Update user name with suffix.
	if count > 0 {
		name = fmt.Sprintf("%s_%d", name, count)
	}

	return name, nil
}

/*
* @Internal
* Generate random nonce for user sign-in.
 */
func generateNonce(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(
			rand.Reader,
			big.NewInt(int64(len(letters))),
		)
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

/*
* @Internal
* Check whether user name is valid.
 */
func isValidUserName(name string) bool {
	// 1. Check whether name consists with only alphabets and under score.
	if b, err := regexp.MatchString(
		`^[A-Za-z]\w+$`,
		name,
	); err != nil || !b {
		return false
	}

	// 2. Check whether another user with same name exists.
	if user := getUserInfoByName(name); user != nil {
		return false
	}

	// 3. Check whether user name contains bad words.
	if isContainingBadWords(name) {
		return false
	}
	return true
}

/*
* @Internal
* Check whether name contains bad words.
 */
func isContainingBadWords(name string) bool {
	for _, word := range config.BAD_WORDS {
		contains := strings.Contains(
			strings.ToUpper(
				strings.ReplaceAll(
					name,
					"_",
					"",
				),
			),
			strings.ToUpper(word),
		)
		if contains {
			return true
		}
	}
	return false
}
