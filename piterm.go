package main
/* Simple Terminal Application for Raspberry Pi
 * http://github.com/coreyshuman/piterm
 * (C) 2016 Corey Shuman
 * 9/6/16
 *
 * License: MIT
 *
 * Usage: piterm serialdev baudrate [(hres)x(vres)]
 *
 * Example: piterm /dev/tty0 9600
 * Example: piterm /dev/tty1 115200 320x240
 */
 
import (
    "os"
    "fmt"
    "sync"
    "strings"
    "strconv"
    "runtime"
	"time"
	"encoding/hex"
    "github.com/mattn/go-gtk/gtk"
	"github.com/coreyshuman/xbeeapi"
)

const timeout = 5
// bb-8 head address
const headAddress = []byte{0x00, 0x13, 0xa2, 0x00, 0x40, 0x90, 0x2a, 0x21}
// bb-8 body address
const bodyAddress = []byte{0x00, 0x13, 0xa2, 0x00, 0x40, 0x90, 0x29, 0x23}

// buffers
var bufAscii *gtk.EntryBuffer = nil
var bufHex *gtk.EntryBuffer = nil

func main() {
	var wg sync.WaitGroup
	quit := make(chan bool)
    hres := 480
    vres := 280
    var start, end gtk.TextIter
	var err error

    fmt.Println("Cores: " + strconv.Itoa(runtime.NumCPU()))
    runtime.GOMAXPROCS(runtime.NumCPU())
    
    if len(os.Args) < 3 {
        fmt.Println("Usage: piterm /dev/tty1 115200 480x240")
        return
    }
    
    dev := os.Args[1]
	baud := os.Args[2]
	baudn, _ := strconv.Atoi(baud)
    
    if(baudn < 1) {
        fmt.Println("Invalid Baud Rate")
        fmt.Println("Usage: piterm /dev/tty1 115200 480x240")
        return
    }
    
    if len(os.Args) > 3 {
        res := strings.Split(os.Args[3], "x")
        if len(res) != 2 {
            fmt.Println("Invalid Resolution Format")
            fmt.Println("Usage: piterm /dev/tty1 115200 480x240")
            return
        }
        hres, _ := strconv.Atoi(res[0])
        vres, _ := strconv.Atoi(res[1])
        if hres < 100 || vres < 100 {
            fmt.Println("Minimum Resolution Must Be 100x100")
            fmt.Println("Usage: piterm /dev/tty1 115200 480x240")
            return
        }
    }
	
	_, err = xbeeapi.Init(devx, baudnx, 1)
	if(err != nil) {
		fmt.Println("Error: " + err.Error())
		return
	}
	// configure xbee api and start job
	xbeeapi.SetupErrorHandler(errorCallback)
	xbeeapi.SetupModemStatusCallback(modemStatusCallback)
	xbeeapi.Begin()
	fmt.Println("XBEE: " + fmt.Sprintf("%d",serialXBEE))
    
    gtk.Init(nil)
    
    // Initialize GUI
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("PiTerm")
	window.SetIconName("gtk-dialog-info")
	window.Connect("destroy", func() {
		quit <- true
		xbeeapi.Close()
		wg.Wait()
		time.Sleep(time.Millisecond*30)
		gtk.MainQuit()
	})
    window.SetSizeRequest(hres, vres)
    vbox := gtk.NewVBox(false, 1)
    hbox1 := gtk.NewHBox(false, 1)
    hbox2 := gtk.NewHBox(false, 1)
    // textbox 1 (ascii)
    swin1 := gtk.NewScrolledWindow(nil, nil)
	swin1.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin1.SetShadowType(gtk.SHADOW_IN)
	textview1 := gtk.NewTextView()
	bufAscii = textview1.GetBuffer()
	bufAscii.GetStartIter(&start)
	bufAscii.GetEndIter(&end)
	//bufAscii.Insert(&end, "Hello")
	swin1.Add(textview1)
	hbox1.Add(swin1)
    // textbox 2 (hex)
    swin2 := gtk.NewScrolledWindow(nil, nil)
	swin2.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin2.SetShadowType(gtk.SHADOW_IN)
	textview2 := gtk.NewTextView()
	bufHex = textview2.GetBuffer()
	bufHex.GetStartIter(&start)
	bufHex.GetEndIter(&end)
	//bufHex.Insert(&end, "World!")
	swin2.Add(textview2)
	hbox1.Add(swin2)
    // textbox and buttons
    textview3 := gtk.NewTextView()
	bufSend := textview3.GetBuffer()
    hbox2.Add(textview3)
    btnSend := gtk.NewButtonWithLabel("Send")
	btnSend.Clicked(func() {
		bufSend.GetStartIter(&start)
		bufSend.GetEndIter(&end)
		sendData := bufSend.GetText(&start, &end, true)
		_, _, err = xbeeapi.SendPacket(headAddress, []byte{0x00, 0x00}, 0x00, []byte(sendData))
		if(err != nil) {
			fmt.Println("Send Error: " + err.Error())
		}
	})
    btnClear := gtk.NewButtonWithLabel("Clear")
	btnClear.Clicked(func() {
		bufAscii.GetStartIter(&start)
		bufAscii.GetEndIter(&end)
		bufAscii.Delete(&start, &end)
		bufHex.GetStartIter(&start)
		bufHex.GetEndIter(&end)
		bufHex.Delete(&start, &end)
	})
    hbox2.Add(btnSend)
    hbox2.Add(btnClear)
    
    vbox.Add(hbox1)
    vbox.Add(hbox2)
    window.Add(vbox)
	window.ShowAll()
    
    go func() {
		wg.Add(1)
		for {
			select {
			case <- quit:
				closeApp()
				wg.Done()
				return
			default:
				time.Sleep(time.Millisecond*30)
			}		
		}
	}()
	
	gtk.Main()
}

func closeApp() {
    
}


/************** Callback Functions ****************/
func errorCallback(e error) {
	fmt.Println(e.Error())
}

var atCommandCallback xbeeapi.ATCommandCallbackFunc = func(frameId byte, data []byte) {
	fmt.Println("AT Response: ")
	fmt.Println(hex.Dump(data))
}

var receivePacketCallback xbeeapi.ReceivePacketCallbackFunc = func(destinationAddress64 [8]byte, destinationAddress16 [2]byte, receiveOptions byte, data []byte) {
	var e gtk.TextIter
	
	bufAscii.GetEndIter(&e)
	bufAscii.Insert(&e, string(data[:]))
	bufHex.GetEndIter(&e)
	bufHex.Insert(&e, hex.EncodeToString(data[:]))
}

var modemStatusCallback xbeeapi.ModemStatusCallbackFunc = func(status byte) {
	modemStatus := xbeeapi.GetModemStatusDescription(status)
	fmt.Println("Modem Status: " + modemStatus)
}