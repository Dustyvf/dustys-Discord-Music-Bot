package main

// bot link https://discordapp.com/api/oauth2/authorize?client_id=516941386369597441&scope=bot&permissions=518208

import (
	"fmt"
	"flag"
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

type lowconf struct {
	Token		string
	DBGLvl		string
	Prefix		[]string
	Name		string
	Database struct {
		DBAddress	string
		DBPass		string
		}
	}

var (
	confName	string
	confToken	string
	confDebug	int
	confPrefix	[]string

	confdbURL	string
	confdbPass	string
	confdbType	string
)

type error interface {
    Error() string
}

const version = "v0.0.0.1:alpha"
const appname = "Dustys Wip Discord Bot"

var useTUI bool
var chk1 int
var messagechk string

type cmddata struct {
	cmdFunc		func(args []string, s disgord.Session, d *disgord.MessageCreate)error
	cmdName		string
	cmdCalls	[]string
	cmdMinDesc	string
	cmdFullDesc string
	cmdFirstChr string
}

var cmdarray = make([]cmddata, 0)

func init() {
	flag.BoolVar(&useTUI, "tui", false, "Use Tui, true/false")
	flag.Parse()
	confDebug = 5
	dlogger.LogOld(0,15,"tui flag set to", strconv.FormatBool(useTUI))

	setupConf()

	dlogger.LogOld(0,99,"Starting up", confName)
	dlogger.LogOld(1,99,"Version", version)
	setupConf();
	dlogger.LogOld(0,15,"Prefix is", confPrefix[0])
}

func prefixCheck(data string) (bool, string) {

	arraylen := len(confPrefix)
	dlogger.LogOld(0,5,"Prefix Amount", strconv.Itoa(arraylen))
	dlogger.LogOld(0,5,"Prefix", fmt.Sprint(confPrefix))

	for i := 0; i < arraylen; i++ {
		if strings.HasPrefix(data, confPrefix[i]) {
			return true, confPrefix[i]
			break
		}
	}
	return false, ""
}

func messageDo(session disgord.Session, data *disgord.MessageCreate) /*(string, string, error)*/ {
	//var responce/*, meta*/ string
	//var err error

	messagechk = data.Message.ID.String()

	ckprfx, prefix := prefixCheck(data.Message.Content)

	if ckprfx {
		msg := strings.Replace(data.Message.Content, prefix, "", -1)
		arg := strings.Fields(msg)
		cmd := arg[0]

			// some debug code
		if confDebug < 5 {
			dlogger.LogOld(0,5,"cmd", cmd)
			dlogger.LogOld(0,5,"args", fmt.Sprint(arg))

			arraylen3 := len(arg)
			dlogger.LogOld(0,5,"Argamt", strconv.Itoa(arraylen3))

			for i := 0; i < arraylen3; i++ {
				dlogger.LogOld(0,5,"Arg", arg[i])
			}
		}

		dlogger.LogOld(0,5,"cmd data",data.Message.Content)

		arraylen := len(cmdarray)
		dlogger.LogOld(0,5,"cmds count", strconv.Itoa(arraylen))

		for i := 0; i < arraylen; i++ {
			if strings.HasPrefix(cmd, cmdarray[i].cmdFirstChr) {

				arraylen2 := len(cmdarray[i].cmdCalls)

				dlogger.LogOld(0,5,"call count", strconv.Itoa(arraylen2))

				for i2 := 0; i2 < arraylen2; i2++ {
				cmdc := cmdarray[i].cmdCalls[i2]
				dlogger.LogOld(0,5,"cmdc", cmdc)
					if cmdc == "" {
						break
					}
					if cmd == cmdc {
						dlogger.LogOld(0,5,"cmd data",cmdc)
						dlogger.LogOld(0,5,cmdc, cmd)
						cmdarray[i].cmdFunc([]string{}, session, data)
					}
				}
			} else {
				if cmdarray[i].cmdFirstChr == "" {
					break
				}
			}
			//dlogger.LogOld(0,5,"cmdchk", cmd)
			//if strings.HasPrefix(msg, cmd) {
			//	dta := strings.Replace(msg, cmd, "", -1)
			//	dlogger.LogOld(0,5,"command", cmd)
			//	dlogger.LogOld(0,5,"arguments", dta)
			//	go cmdcorehandler(cmd, dta , session, data)
			//	break
			//}
		}

		//responce = "hello"
		//msg.RespondString(session, responce)
	}

	//return responce, meta, err
}

func main() {

	session, err := disgord.NewSession(&disgord.Config{
		Token: confToken,
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

		user, err := session.GetCurrentUser()
		if err != nil {
			dlogger.LogOld(30,25,"Error getting current user","")
		}
		fmt.Println(user.ID)
		fmt.Println(data.Message.Author)
		if data.Message.Author.ID != user.ID {
			if (msg.ID.String() != messagechk) {
				go messageDo(session, data)
				}
			}
	})

	err = session.Connect()
	if err != nil {
		dlogger.LogOld(50, 999999, "Discord Session error", "")
		dlogger.LogOld(51, 999999, err.Error(), "")
		panic(err)
	}

	dlogger.SetLevels(confDebug)
	tst := dlogger.Check()
	dlogger.LogOld(0,15, "debug check", strconv.Itoa(tst))
	dlogger.LogExtraInfo(15,"test","")

	dlogger.LogOld(0,15, "Running under user", myself.String())

	session.DisconnectOnInterrupt()
}
