package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/gorilla/mux"
	"github.com/pkg/browser"
)

const (
	startedByOS = "__byOS" // Started by OS flag
	appName     = "LocalNetTools"
)

var (
	port               string
	shuttingDown       bool
	debugging          bool
	startedWithOS      bool
	startWithOsFlagged bool
	timeElapsed        int64
	remaningTimeValue  int64
)

// --------------------- ROUTES ---------------------- //

// Start page. TODO: GUI
func index(w http.ResponseWriter, r *http.Request) {
	// message := fmt.Sprintf("%s v0.2.0", appName)
	// newResponse(w, respAPI{Message: message})
	b := MustAsset("webui/dist/index.html")
	// if err != nil {
	// 	panic(err)
	// }

	w.Write(b)

	go _debug()
}

func getUserDetail(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	user, err := user.Current()
	if err != nil {
		newResponse(w, respAPI{Message: "Error getting user: " + err.Error()})
		return
	}

	newResponse(w, respAPI{Message: user.Name})
	go _debug()
}

// Shutdown PC
func shutDown(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	// Try to get delay to shutdown PC
	message := r.URL.Path
	paths := strings.Split(message, "/")
	message = paths[len(paths)-1]

	var cmd *exec.Cmd

	delayTime, err := strconv.ParseInt(message, 10, 64)
	oriDelay := delayTime
	// User not setted a valid number to delay shutdown
	if err != nil {
		_setShutdown(59)
		if runtime.GOOS == "windows" {
			cmd = exec.Command("shutdown", "/s", "/f")
		} else {
			cmd = exec.Command("/sbin/shutdown", "-h", "1")
		}
		err = cmd.Run()
		if err != nil {
			newResponse(w, respAPI{Message: "Error: " + err.Error(), Error: true})
			return
		}
		newResponse(w, respAPI{Message: "Shutting down PC after one minute"})
		return

	} else {
		// User setted a valid number to delay shutdown
		inMinutes := delayTime * 60
		_setShutdown(oriDelay)
		respMessage := fmt.Sprintf("%v", oriDelay)
		if runtime.GOOS == "windows" {
			cmd = exec.Command("shutdown", "/s", "/f", "/t", fmt.Sprintf("%v", inMinutes))
		} else {
			cmd = exec.Command("/sbin/shutdown", "-h", fmt.Sprintf("%v", delayTime))
		}

		err := cmd.Run()

		if err != nil {
			newResponse(w, respAPI{Message: "Error: " + err.Error(), Error: true})
			return
		}

		newResponse(w, respAPI{Message: respMessage})
		return
	}

	if err != nil {
		newResponse(w, respAPI{Message: "Error: " + err.Error(), Error: true})
		return
	}

	newResponse(w, respAPI{Message: "Unknown response", Error: true})

	go _debug()
}

// Cancel shutdown
func shutdownCancel(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if !shuttingDown {
		newResponse(w, respAPI{Message: "No shutdown scheduled", Error: true})
		go _debug()
		return
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("shutdown", "/a")
	} else {
		cmd = exec.Command("/sbin/shutdown", "-c")
	}

	cmd.Run()
	remaningTimeValue = 0
	shuttingDown = false
	newResponse(w, respAPI{Message: "Shutdown was canceled"})
	go _debug()
}

// Exit program
func exit(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	newResponse(w, respAPI{Message: fmt.Sprintf("%s exited", appName)})
	go func() {
		time.Sleep(time.Second * 2)
		os.Exit(0)
	}()
}

// Return time since user was logged in
func sinceLaunch(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if startedWithOS {
		newResponse(w, respAPI{Message: fmt.Sprintf("%v", timeElapsed)})
	} else {
		newResponse(w, respAPI{Message: "No data to show. Enable option to launch with O.S. with command line: -s"})
	}
}

// Recieve time left for shutdown the O.S.
func shutdownTimeLeft(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	remaningTime := remaningTimeValue - time.Now().Unix()
	if remaningTimeValue != 0 {
		newResponse(w, respAPI{Message: fmt.Sprintf("%v", remaningTime)})
	} else {
		newResponse(w, respAPI{Message: "No shutdown scheduled", Error: true})
	}
}

// --------------------- END ROUTES ---------------------- //

// ------------------ ROUTERS HELPERS -------------------- //

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// Set remaning time to shwtdown
func _setShutdown(secondsToShutdown int64) {
	shuttingDown = true
	remaningTimeValue = time.Now().Unix() + (secondsToShutdown * 60)
}

// Calculate the time elapsed since O.S. started
func _elapsedTime(found bool, result string) {
	if found {
		go func() {
			for {
				time.Sleep(time.Second)
				timeElapsed++
			}
		}()
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	headers := w.Header()
	if strings.HasSuffix(filePath, "css") {
		headers["Content-Type"] = []string{"text/css"}
	} else if strings.HasSuffix(filePath, "js") {
		headers["Content-Type"] = []string{"application/javascript"}
	} else if strings.HasSuffix(filePath, "woff") {
		headers["Content-Type"] = []string{"application/x-font-woff"}
	} else if strings.HasSuffix(filePath, "woff2") {
		headers["Content-Type"] = []string{"application/x-font-woff2"}
	}

	// b := MustAsset("public/static/app.css")
	println(filePath)
	b := MustAsset("webui/dist" + filePath)
	w.Write(b)
}

// ------------------ ROUTERS HELPERS -------------------- //

// Debug
func _debug() {
	if !debugging {
		return
	}
	time.Sleep(time.Second * 2)
	cmd := exec.Command("shutdown", "/a")
	cmd.Run()

	os.Exit(0)
}

// Check if has 'start with O.S.' and register
func _runAtStartUpRegister(found bool, result string) {
	if found {
		_addToStartup()
	}
}

// Handle custom port number
func _customPortHandler(found bool, result string) {
	if found {
		newPort, err := strconv.Atoi(result)
		// If it isnt valid number use default port number
		if err != nil {
			return
		}

		if newPort < 1024 || newPort > 65535 {
			return
		}
		port = fmt.Sprintf(":%v", newPort)
	}
}

// Setup handle events for systray icon
func _setupSystray() {
	systray.SetIcon(MustAsset("neticon.ico"))
	systray.SetTooltip("LocalNetTools")
	openUI := systray.AddMenuItem("Open UI", "Open UI in default browser")
	startUp := systray.AddMenuItem("Run in startup", "Run on system startup")
	systray.AddSeparator()
	quit := systray.AddMenuItem("Quit", "Quit the app")

	if startedWithOS || startWithOsFlagged {
		startUp.Check()
	}

	go func() {
		for {
			select {
			case <-openUI.ClickedCh:
				conn, err := net.Dial("udp", "8.8.8.8:80")
				if err != nil {
					log.Fatal(err)
				}
				defer conn.Close()

				localAddr := conn.LocalAddr().(*net.UDPAddr)

				address := localAddr.IP.String()
				browser.OpenURL(fmt.Sprintf("http://%s:1302/", address))
			case <-quit.ClickedCh:
				_exitSystray()
			case <-startUp.ClickedCh:
				if startUp.Checked() {
					startUp.Uncheck()
					_removeFromStartup()
				} else {
					startUp.Check()
					_addToStartup()
				}
			}
		}
	}()
}

// Exit program
func _exitSystray() {
	systray.Quit()
	os.Exit(0)
}

func main() {

	_startedWithOS()

	go func() {
		// System icon
		systray.Run(_setupSystray, _exitSystray)
	}()

	// Define default port
	port = "1302"

	// Test for debug flag at start [help development only]
	checkArg("__debug-development-only__", func(found bool, result string) {
		if found {
			debugging = true
		}
	})

	// Check if has START WITH OS arg
	if !checkArg(startedByOS, _elapsedTime) {
		checkArg("-s", _runAtStartUpRegister)
		checkArg("-p=", _customPortHandler)
	}

	router := mux.NewRouter().StrictSlash(true)

	router.PathPrefix("/js/").HandlerFunc(staticHandler)
	router.PathPrefix("/css/").HandlerFunc(staticHandler)
	router.PathPrefix("/fonts/").HandlerFunc(staticHandler)

	// Handles api pathsgetRemaningTime
	router.HandleFunc("/api/exit", exit)
	router.HandleFunc("/api/sincelaunch", sinceLaunch)
	router.HandleFunc("/api/shutdown/timeleft", shutdownTimeLeft)
	router.HandleFunc("/api/shutdown/cancel", shutdownCancel)
	router.HandleFunc("/api/shutdown/a", shutdownCancel)
	router.HandleFunc("/api/shutdown/c", shutdownCancel)
	router.HandleFunc("/api/shutdown/{delay}", shutDown)
	router.HandleFunc("/api/getuserdetail/", getUserDetail)

	router.HandleFunc("/", index)

	// Show info about server port
	fmt.Printf("Running on port %s...", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		panic(err)
	}

}
