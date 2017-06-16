package main

import (
	"flag"
	"log"
	"time"
	"fmt"
	"encoding/json"
	"os"

	"github.com/nats-io/go-nats"
)


const CLR_0 = "\x1b[30;1m"
const CLR_R = "\x1b[31;1m"
const CLR_G = "\033[32m"
const CLR_Y = "\x1b[33;1m"
const CLR_B = "\x1b[34;1m"
const CLR_M = "\x1b[35;1m"
const CLR_C = "\x1b[36;1m"
const CLR_W = "\x1b[37;1m"
const CLR_N = "\033[0m"

var (
	replyEvent chan *nats.Msg
	urls = flag.String("s", "", "")
        user = flag.String("u", "nats" , "")
        pass = flag.String("p", "", "")
        subj = flag.String("sub", "", "")
        message = flag.String("m", "", "")
        timeout = flag.Int("t", 10, "")
)

// NOTE: Use tls scheme for TLS, e.g. nats-req -s tls://demo.nats.io:4443 foo hello
func usage() {
	log.Fatalf(`
Usage:

./natsping -s [ ip:port | ip:port,ip:port ] -u username -p password -sub subject -m message -t <timeout value in seconds default=10>

Use Case 1 
	bosh director is reporting agent ping timeout after <X> seconds even though the agent is up and running on the remove vm. In this case you may want to ping the remote bosh agent manualy
	Get the agent id from the %[1]sbosh vms --details%[2]s you want ping.  The agent id will be used in the request subject as %[1]sagent.<agent id>%[2]s.  If your agent id is %[1]s05e3468d-72e1-4796-a871-0c143c25013a%[2]s then you could send a ping with the following command.

%[3]s./natsping -s "10.193.67.11:4222" -u nats -p "password" -sub "agent.05e3468d-72e1-4796-a871-0c143c25013a" -m '{"method":"ping","arguments":[], "reply_to": "agent.reply_to_natsping"}'%[2]s
`, CLR_G, CLR_N, CLR_M)
}

func printMsg(m *nats.Msg) {
	log.Printf("[%s]Received on subject [%s]: '%s'\n", time.Now(), m.Subject, string(m.Data))
	if m.Subject == getReplyTo(*message) {
		replyEvent <- m
	}
}

func getReplyTo(m string) string {
	type Data struct {
		ReplyTo string `json:"reply_to"`
	}
	d := Data{}
	err := json.Unmarshal([]byte(m), &d)
	if err != nil {
		log.Fatalf("Failed to unmarshal message [%s]: %s\n", m, err)
	}
	if d.ReplyTo == "" {
		log.Fatalf("reply_to not set in message")
	}
	return d.ReplyTo
}

func startWatchReply() {
	select {
		case <-replyEvent:
			log.Println("Reply received successfully!")
			os.Exit(0)
	}
}

func main() {
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if *urls == "" || *user == "" || *pass == "" || *subj == "" || *message == "" {
		usage()
	}
	secureUrl := fmt.Sprintf("nats://%s:%s@%s", *user, *pass, *urls)
	fmt.Printf("Using url %s\n", secureUrl)

	nc, err := nats.Connect(secureUrl)
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}
	defer nc.Close()

	nc.Subscribe(*subj, printMsg)
	nc.Subscribe(getReplyTo(*message), printMsg)

	/*nc.Subscribe(*subj, func(msg *nats.Msg) {
		i += 1
		printMsg(msg,i)
	})
	
	y := 0
        nc.Subscribe("agent.reply_to_natsping", func(msg *nats.Msg) {
                y += 1
                printMsg(msg,y)
        }) */

	
	replyEvent = make(chan *nats.Msg, 0)
	go startWatchReply()

	log.Printf("Published [%s] : '%s'\n", *subj, *message)
	msg, err := nc.Request(*subj, []byte(*message), time.Duration(*timeout)*time.Second)
	if err != nil {
		if nc.LastError() != nil {
			log.Fatalf("Error in Request: %v\n", nc.LastError())
		}
		log.Printf("Error in Request: %v\n", err)
	} else {
	  log.Printf("Received [%v] : '%s'\n", msg.Subject, string(msg.Data))
	}
}
