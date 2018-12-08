package main

import (
	"flag"
	"fmt"
	//"os"
	//"os/signal"
	//"syscall"
	//"io/ioutil"
	//"encoding/json"
	"strings"
	//"time"
	"strconv"
	//"reflect"
	"./dlogger"

	"github.com/andersfylling/disgord"
)

type Lowconf struct {
	Token  string
	DBGLvl string
	Prefix []string
	Name   string
}

var (
	conf_Name   string
	conf_Token  string
	conf_Debug  int
	conf_Prefix []string
)

const version = "v0.0.0.1:alpha"
const appname = "Wip Discord Bot"

var UseTUI bool // not currently used
var chk1 int = 0

var messagechk1 string = "~~~~~~" // thing for a thing to prevent duplicate message respoces (may be uneeded, not sure yet)
var messagechk2 string = "~~~~~~~"

func init() {
	flag.BoolVar(&UseTUI, "tui", false, "Use Tui, true/false")
	flag.Parse()
	conf_Debug = 5
	dlogger.LogOld(0, 15, "tui flag set to", strconv.FormatBool(UseTUI))

	setupConf()

	dlogger.LogOld(0, 99, "Starting up", conf_Name)
	dlogger.LogOld(1, 99, "Version", version)
	setupConf()
	dlogger.LogOld(0, 15, "Prefix is", conf_Prefix[0])
}

func prefixCheck(data string) (bool, string) {
	prearraylen := len(conf_Prefix)
	dlogger.LogOld(0, 5, "Prefix Amount", strconv.Itoa(prearraylen))
	for i := 0; i < prearraylen; i++ {
		pfx := conf_Prefix[i]
		dlogger.LogOld(0, 5, "Prefix", pfx)
		if strings.HasPrefix(data, pfx) {
			return true, pfx
			break
		}
	}
	return false, ""
}

func messageDo(message string, session disgord.Session, data *disgord.MessageCreate) {
	msg := data.Message

	messagechk1 = msg.Content

	ckprfx, prefix := prefixCheck(message)

	if ckprfx {
		message2 := strings.Replace(message, prefix, "", -1)
		dlogger.LogOld(0, 5, "cmd data", message2)

		prearraylen := len(corecmdslist)
		dlogger.LogOld(0, 5, "core cmds count", strconv.Itoa(prearraylen))
		for i := 0; i < prearraylen; i++ {
			cmd := corecmdslist[i]
			if cmd == "" {
				break
			}
			dlogger.LogOld(0, 5, "cmdchk", cmd)
			if strings.HasPrefix(message2, cmd) {
				dta := strings.Replace(message2, cmd, "", -1)
				dlogger.LogOld(0, 5, "command", cmd)
				dlogger.LogOld(0, 5, "arguments", dta)
				go cmdcorehandler(cmd, dta, session, data)
				break
			}
		}

	}

	messagechk1 = "~~~~~~"
}

func main() {

	session, err := disgord.NewSession(&disgord.Config{
		Token: conf_Token,
		Debug: true,
	})
	if err != nil {
		dlogger.LogOld(50, 999999, "Failed to open discord session", "")
		dlogger.LogOld(51, 999999, err.Error(), "")
		panic(err)
	}

	myself, err := session.GetCurrentUser()
	if err != nil {
		dlogger.LogOld(50, 999999, "Discord Session error", "")
		dlogger.LogOld(51, 999999, err.Error(), "")
		panic(err)
	}

	session.On(disgord.EventMessageCreate, func(session disgord.Session, data *disgord.MessageCreate) {
		msg := data.Message
		dlogger.LogOld(5, 15, "Message recived", msg.Content)

		messagechk2 = msg.Content

		user, err := session.GetCurrentUser()
		if err != nil {
			dlogger.LogOld(30, 25, "Error getting current user", "")
		}
		fmt.Println(user.ID)
		fmt.Println(data.Message.Author)
		if data.Message.Author.ID != user.ID {
			if messagechk1 != messagechk2 {
				go messageDo(msg.Content, session, data)
			}
		}
	})

	err = session.Connect()
	if err != nil {
		dlogger.LogOld(50, 999999, "Discord Session error", "")
		dlogger.LogOld(51, 999999, err.Error(), "")
		panic(err)
	}

	dlogger.SetLevels(conf_Debug)
	tst := dlogger.Check()
	dlogger.LogOld(0, 15, "debug check", strconv.Itoa(tst))
	dlogger.LogExtraInfo(15, "test", "")

	dlogger.LogOld(0, 15, "Running under user", myself.String())

	session.DisconnectOnInterrupt()
}
