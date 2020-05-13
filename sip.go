// Copyright (c) 2018 Vasily Suvorov, http://bazil.pro <gbazil@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
//

package sipd

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strings"
)

const (
	newLine   = "\r\n"
	seperator = ": "
)

var (
	methodFromRegex = regexp.MustCompile(`\S*`)
	nameFromRegex   = regexp.MustCompile(`sip:[^@]*`)
	addrFromRegex1  = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+:\d+`)
	addrFromRegex2  = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	addrFromRegex3  = regexp.MustCompile(`@[\w.-]*`)
	valueFromRegex  = regexp.MustCompile(`sip:[^@]*`)
)

// Message is the basic SIP Message
type Message map[string]string

// Parse populates the message using the string
func (m Message) Parse(s string) error {
	var f bool
	for i, v := range strings.Split(s, newLine) {
		if i == 0 {
			m["Request"] = v
			continue
		}
		if len(v) <= 0 {
			f = true
			continue
		}
		if f {
			m["Content"] += v + newLine
			continue
		}
		pair := strings.Split(v, seperator)
		if len(pair) < 2 {
			return fmt.Errorf("unexpected line: %q", v)
		}
		if len(m[pair[0]]) > 0 {
			m[pair[0]] = m[pair[0]] + "," + pair[1]
		} else {
			m[pair[0]] = pair[1]
		}
	}
	return nil
}

func (m Message) String() string {
	var s string
	for k, v := range m {
		if k == "Via" {
			for _, vv := range strings.Split(v, ",") {
				s += k + seperator + vv + newLine
			}
		} else if k != "Request" && k != "Content" {
			s += k + seperator + v + newLine
		}
	}

	return m["Request"] + newLine + s + newLine + m["Content"]
}

func (m Message) Clone() (clone Message) {
	clone = make(Message)
	for key, val := range m {
		clone[key] = val
	}
	return
}

func (m Message) MethodFrom(key string) string {
	return methodFromRegex.FindString(m[key])
}

func (m Message) NameFrom(key string) string {
	return strings.TrimPrefix(nameFromRegex.FindString(m[key]), "sip:")
}

func (m Message) NameFor(key, name string) {
	if oldname := m.NameFrom(key); oldname != "" {
		m[key] = strings.Replace(m[key], "sip:"+oldname, "sip:"+name, 1)
	}
}

func (m Message) AddrFrom(key string) (addr string) {
	if addr = addrFromRegex1.FindString(m[key]); addr != "" {
		return
	}
	if addr = addrFromRegex2.FindString(m[key]); addr != "" {
		return
	}
	addr = strings.TrimPrefix(addrFromRegex3.FindString(m[key]), "@")
	return
}

func (m Message) AddrFor(key, addr string) {
	if s := m.AddrFrom(key); s != "" {
		m[key] = strings.Replace(m[key], s, addr, 1)
	}
}

func (m Message) ValueFrom(key, value string) string { // slow
	s := strings.Replace(m[key], `"`, ``, -1)
	return strings.TrimPrefix(regexp.MustCompile(value+`=[^ ,;>]*`).FindString(s), value+`=`)
}

// PrepareReply only updates the message for the reply
func (m Message) PrepareReply() {
	delete(m, "Allow")
	delete(m, "Supported")

	delete(m, "Content")
	delete(m, "Content-Type")
	m["Content-Length"] = "0"
	return
}

func (m Message) Digest(secret string) string {
	// HA1=MD5(username:realm:password) HA2=MD5(method:digestURI) response=MD5(HA1:nonce:HA2)
	b1 := []byte(m.ValueFrom("Authorization", "username") + ":" + m.ValueFrom("Authorization", "realm") + ":" + secret)
	h1 := fmt.Sprintf("%x", md5.Sum(b1))

	b2 := []byte(m.MethodFrom("Request") + ":" + m.ValueFrom("Authorization", "uri"))
	h2 := fmt.Sprintf("%x", md5.Sum(b2))

	b3 := []byte(h1 + ":" + m.ValueFrom("Authorization", "nonce") + ":" + h2)
	return fmt.Sprintf("%x", md5.Sum(b3))
}

// MethodFrom will return the SIP method form the given string
func MethodFrom(value string) string {
	return methodFromRegex.FindString(value)
}

// NameFrom will return the SIP name form the given string
func NameFrom(value string) string {
	return strings.TrimPrefix(nameFromRegex.FindString(value), "sip:")
}

// HostFrom will return the host form the given string
func HostFrom(value string) (addr string) {
	if addr = addrFromRegex1.FindString(value); addr != "" {
		return
	}
	if addr = addrFromRegex2.FindString(value); addr != "" {
		return
	}
	addr = strings.TrimPrefix(addrFromRegex3.FindString(value), "@")
	return
}
