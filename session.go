package odm

import (
	"crypto/tls"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
)

type conInfo struct {
	Host     string
	User     string
	Pass     string
	Protocol string
}

func dial(uri string) (*mgo.Session, error) {
	dialInfo, err := prepareInfo(uri)
	if err != nil {
		return nil, err
	}

	return mgo.DialWithInfo(dialInfo)
}

func prepareInfo(uri string) (*mgo.DialInfo, error) {
	info, err := parseURI(uri)
	if err != nil {
		return nil, err
	}

	dialInfo := &mgo.DialInfo{
		Password: info.Pass,
		Username: info.User,
		FailFast: true,
		Timeout:  10 * time.Second}

	if info.Protocol == "mongodb+srv" {
		hosts, err := dnsLookup(info.Host)
		if err != nil {
			return nil, err
		}
		dialInfo.Addrs = hosts
	} else {
		dialInfo.Addrs = parseHosts(info.Host)
	}

	tlsConfig := &tls.Config{}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		return tls.Dial("tcp", addr.String(), tlsConfig)
	}

	return dialInfo, err
}

func parseHosts(uri string) []string {
	hosts := strings.Split(uri, ",")
	return hosts
}

func dnsLookup(host string) ([]string, error) {
	_, addr, err := net.LookupSRV("mongodb", "tcp", host)
	if err != nil {
		return nil, err
	}

	hosts := []string{}
	for _, a := range addr {
		hostname := a.Target[0 : len(a.Target)-1]
		hosts = append(hosts, fmt.Sprintf("%s:%d", hostname, a.Port))
	}

	return hosts, nil
}

func parseURI(uri string) (*conInfo, error) {
	reg, err := regexp.Compile("(.+):\\/\\/(?:(.+)@)?([a-zA-Z0-9\\-\\.\\:\\,\\_]+)")
	if err != nil {
		panic(err)
	}

	matches := reg.FindStringSubmatch(uri)
	var ret *conInfo

	if len(matches) < 3 || len(matches) > 5 {
		return nil, fmt.Errorf("Invalid uri")
	}

	ret = &conInfo{Protocol: matches[1]}
	if len(matches) == 3 {
		ret.Host = matches[2]
		return ret, nil
	}

	userPass := strings.Split(matches[2], ":")
	if len(userPass) == 2 {
		ret.User = userPass[0]
		ret.Pass = userPass[1]
	}

	ret.Host = matches[3]

	return ret, nil
}
