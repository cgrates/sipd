package sipingo

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestNewMessage(t *testing.T) {
	inviteMessage := "INVITE sip:1002@192.168.58.203 SIP/2.0\r\nCall-ID: 4d4d84b0cc83fc90aca41e295cd8ff43@0:0:0:0:0:0:0:0\r\nCSeq: 2 INVITE\r\nFrom: \"1001\" <sip:1001@192.168.58.203>;tag=99f35805\r\nTo: <sip:1002@192.168.58.203>\r\nMax-Forwards: 70\r\nContact: \"1001\" <sip:1001@192.168.58.201:5060;transport=udp;registering_acc=192_168_58_203>\r\nUser-Agent: Jitsi2.11.20200408Linux\r\nContent-Type: application/sdp\r\nVia: SIP/2.0/UDP 192.168.58.201:5060;branch=z9hG4bK-393139-939e89686023b86822cb942ede452b62\r\nProxy-Authorization: Digest username=\"1001\",realm=\"192.168.58.203\",nonce=\"XruO2167ja8uRODnSv8aXqv+/hqPJiXh\",uri=\"sip:1002@192.168.58.203\",response=\"5b814c709d1541d72ea778599c2e48a4\"\r\nContent-Length: 897\r\n\r\nv=0\r\no=1001-jitsi.org 0 0 IN IP4 192.168.58.201\r\ns=-\r\nc=IN IP4 192.168.58.201\r\nt=0 0\r\nm=audio 5000 RTP/AVP 96 97 98 9 100 102 0 8 103 3 104 101\r\na=rtpmap:96 opus/48000/2\r\na=fmtp:96 usedtx=1\r\na=ptime:20\r\na=rtpmap:97 SILK/24000\r\na=rtpmap:98 SILK/16000\r\na=rtpmap:9 G722/8000\r\na=rtpmap:100 speex/32000\r\na=rtpmap:102 speex/16000\r\na=rtpmap:0 PCMU/8000\r\na=rtpmap:8 PCMA/8000\r\na=rtpmap:103 iLBC/8000\r\na=rtpmap:3 GSM/8000\r\na=rtpmap:104 speex/8000\r\na=rtpmap:101 telephone-event/8000\r\na=extmap:1 urn:ietf:params:rtp-hdrext:csrc-audio-level\r\na=extmap:2 urn:ietf:params:rtp-hdrext:ssrc-audio-level\r\na=rtcp-xr:voip-metrics\r\nm=video 5002 RTP/AVP 105 99\r\na=recvonly\r\na=rtpmap:105 h264/90000\r\na=fmtp:105 profile-level-id=42E01f;packetization-mode=1\r\na=imageattr:105 send * recv [x=[1:1920],y=[1:1080]]\r\na=rtpmap:\r\n"

	message, err := NewMessage(inviteMessage)
	if err != nil {
		t.Fatal(err)
	}
	expected := Message{
		"Request":             "INVITE sip:1002@192.168.58.203 SIP/2.0",
		"From":                `"1001" <sip:1001@192.168.58.203>;tag=99f35805`,
		"Contact":             `"1001" <sip:1001@192.168.58.201:5060;transport=udp;registering_acc=192_168_58_203>`,
		"Proxy-Authorization": `Digest username="1001",realm="192.168.58.203",nonce="XruO2167ja8uRODnSv8aXqv+/hqPJiXh",uri="sip:1002@192.168.58.203",response="5b814c709d1541d72ea778599c2e48a4"`,
		"To":                  `<sip:1002@192.168.58.203>`,
		"Max-Forwards":        "70",
		"User-Agent":          "Jitsi2.11.20200408Linux",
		"Content-Type":        "application/sdp",
		"Via":                 "SIP/2.0/UDP 192.168.58.201:5060;branch=z9hG4bK-393139-939e89686023b86822cb942ede452b62",
		"Call-ID":             "4d4d84b0cc83fc90aca41e295cd8ff43@0:0:0:0:0:0:0:0",
		"CSeq":                "2 INVITE",
		"Content-Length":      "897",
		"Content":             "v=0\r\no=1001-jitsi.org 0 0 IN IP4 192.168.58.201\r\ns=-\r\nc=IN IP4 192.168.58.201\r\nt=0 0\r\nm=audio 5000 RTP/AVP 96 97 98 9 100 102 0 8 103 3 104 101\r\na=rtpmap:96 opus/48000/2\r\na=fmtp:96 usedtx=1\r\na=ptime:20\r\na=rtpmap:97 SILK/24000\r\na=rtpmap:98 SILK/16000\r\na=rtpmap:9 G722/8000\r\na=rtpmap:100 speex/32000\r\na=rtpmap:102 speex/16000\r\na=rtpmap:0 PCMU/8000\r\na=rtpmap:8 PCMA/8000\r\na=rtpmap:103 iLBC/8000\r\na=rtpmap:3 GSM/8000\r\na=rtpmap:104 speex/8000\r\na=rtpmap:101 telephone-event/8000\r\na=extmap:1 urn:ietf:params:rtp-hdrext:csrc-audio-level\r\na=extmap:2 urn:ietf:params:rtp-hdrext:ssrc-audio-level\r\na=rtcp-xr:voip-metrics\r\nm=video 5002 RTP/AVP 105 99\r\na=recvonly\r\na=rtpmap:105 h264/90000\r\na=fmtp:105 profile-level-id=42E01f;packetization-mode=1\r\na=imageattr:105 send * recv [x=[1:1920],y=[1:1080]]\r\na=rtpmap:\r\n",
	}
	if !reflect.DeepEqual(message, expected) {
		t.Errorf("Expected %s, received: %s", expected, message)
	}

	expectedStrSlice := strings.Split(inviteMessage, newLine) // split and sort the lines
	sort.Strings(expectedStrSlice)

	strMessage := message.String()
	strMessageSlice := strings.Split(strMessage, newLine) // split and sort the lines
	sort.Strings(strMessageSlice)
	if !reflect.DeepEqual(expectedStrSlice, strMessageSlice) {
		t.Errorf("Expected %q, received: %q", inviteMessage, strMessage)
	}

	expectedMethod := "INVITE"
	if rply := message.MethodFrom("Request"); !reflect.DeepEqual(expectedMethod, rply) {
		t.Errorf("Expected %s, received: %s", expectedMethod, rply)
	}

	expectedUser := "1001"
	if rply := message.UserFrom("From"); !reflect.DeepEqual(expectedUser, rply) {
		t.Errorf("Expected %s, received: %s", expectedUser, rply)
	}
	expectedHost := "192.168.58.203"
	if rply := message.HostFrom("From"); !reflect.DeepEqual(expectedHost, rply) {
		t.Errorf("Expected %s, received: %s", expectedHost, rply)
	}

	expected = Message{
		"Request":             "INVITE sip:1002@192.168.58.203 SIP/2.0",
		"From":                `"1001" <sip:1001@192.168.58.203>;tag=99f35805`,
		"Contact":             `"1001" <sip:1001@192.168.58.201:5060;transport=udp;registering_acc=192_168_58_203>`,
		"Proxy-Authorization": `Digest username="1001",realm="192.168.58.203",nonce="XruO2167ja8uRODnSv8aXqv+/hqPJiXh",uri="sip:1002@192.168.58.203",response="5b814c709d1541d72ea778599c2e48a4"`,
		"To":                  `<sip:1002@192.168.58.203>`,
		"Max-Forwards":        "70",
		"User-Agent":          "Jitsi2.11.20200408Linux",
		"Via":                 "SIP/2.0/UDP 192.168.58.201:5060;branch=z9hG4bK-393139-939e89686023b86822cb942ede452b62",
		"Call-ID":             "4d4d84b0cc83fc90aca41e295cd8ff43@0:0:0:0:0:0:0:0",
		"CSeq":                "2 INVITE",
		"Content-Length":      "0",
	}
	message.PrepareReply()
	if !reflect.DeepEqual(message, expected) {
		t.Errorf("Expected %s, received: %s", expected, message)
	}
}

func TestNewMessage2(t *testing.T) {
	inviteMessage := "INVITE sip:1002@192.168.58.203 SIP/2.0\r\nCall-ID\r\n"
	expectedErr := `unexpected line: "Call-ID"`
	if _, err := NewMessage(inviteMessage); err == nil || err.Error() != expectedErr {
		t.Errorf("Expected err: %q,received %v", expectedErr, err)
	}

	inviteMessage = "INVITE sip:1002@192.168.58.203 SIP/2.0\r\nContact: \"1001\" <sip:1001@192.168.58.201:5060>\r\nContact: \"1002\" <sip:1002@192.168.58.201:5060>"
	message, err := NewMessage(inviteMessage)
	if err != nil {
		t.Fatal(err)
	}
	expected := Message{
		"Request": "INVITE sip:1002@192.168.58.203 SIP/2.0",
		"Contact": `"1001" <sip:1001@192.168.58.201:5060>,"1002" <sip:1002@192.168.58.201:5060>`,
	}
	if !reflect.DeepEqual(message, expected) {
		t.Errorf("Expected %s, received: %s", expected, message)
	}
	clndMessage := message.Clone()
	if !reflect.DeepEqual(clndMessage, expected) {
		t.Errorf("Expected %s, received: %s", expected, clndMessage)
	}
	message["To"] = `"1001" <sip:1001@192.168.56.203>`
	if !reflect.DeepEqual(clndMessage, expected) {
		t.Errorf("Expected %s, received: %s", expected, clndMessage)
	}
}

func TestHostFrom(t *testing.T) {
	val := "INVITE sip:1002@192.168.58.203 SIP/2.0"
	expected := "192.168.58.203"
	if rply := HostFrom(val); rply != expected {
		t.Errorf("Expected %q, received: %q", rply, expected)
	}

	val = "INVITE sip:1002@192.168.58.203:5060 SIP/2.0"
	expected = "192.168.58.203:5060"
	if rply := HostFrom(val); rply != expected {
		t.Errorf("Expected %q, received: %q", rply, expected)
	}

	val = "INVITE sip:1002@cgrates.org SIP/2.0"
	expected = "cgrates.org"
	if rply := HostFrom(val); rply != expected {
		t.Errorf("Expected %q, received: %q", rply, expected)
	}
}
