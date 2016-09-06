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
    "strings"
    "github.com/coreyshuman/serial"
    "github.com/mattn/go-gtk/gtk"
)

const timeout = 5

func main() {
	var wg sync.WaitGroup
	quit := make(chan bool)
    hres := 480
    vres := 280
    var start, end gtk.TextIter
    
    if len(os.Args) < 3 {
        fmt.Println("Usage: piterm /dev/tty1 115200 480x240")
        return
    }
    
    dev := os.Args[1]
	baud := os.Args[2]
	baudn, _ := strconv.Atoi(baudx)
    
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
    
    serial.Init()
	sid, err := serial.Connect(dev, baudn, timeout)
    if(err != nil) {
		fmt.Println("Serial Connection Failed: " + err.Error())
		return
	}
    
    gtk.Init(nil)
    
    // Initialize GUI
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("PiTerm")
	window.SetIconName("gtk-dialog-info")
	window.Connect("destroy", func() {
		quit <- true
		wg.Wait()
		gtk.MainQuit()
	})
    window.SetSizeRequest(hres, vres)
    vbox := gtk.NewVBox(false, 1)
    hbox1 := gtk.NewHBox(false, 1)
    hbox2 := gtk.NewVBox(false, 1)
    // textbox 1 (ascii)
    swin1 := gtk.NewScrolledWindow(nil, nil)
	swin1.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin1.SetShadowType(gtk.SHADOW_IN)
	textview1 := gtk.NewTextView()
	var start, end gtk.TextIter
	bufAscii := textview1.GetBuffer()
	//buffer.GetStartIter(&start)
	buffer.GetEndIter(&end)
	buffer.Insert(&end, "Hello")
	swin1.Add(textview1)
	hbox1.Add(swin1)
    // textbox 2 (hex)
    swin2 := gtk.NewScrolledWindow(nil, nil)
	swin2.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin2.SetShadowType(gtk.SHADOW_IN)
	textview2 := gtk.NewTextView()
	var start, end gtk.TextIter
	bufHex := textview2.GetBuffer()
	//buffer.GetStartIter(&start)
	buffer.GetEndIter(&end)
	buffer.Insert(&end, "World!")
	swin2.Add(textview2)
	hbox1.Add(swin2)
    // textbox and buttons
    textview3 := gtk.NewTextView()
    hbox2.Add(textview3)
    btnSend := gtk.NewButtonWithLabel("Send")
    btnClear := gtk.NewButtonWithLabel("Clear")
    hbox2.Add(btnSend)
    hbox2.Add(btnClear)
    
    vbox.Add(hbox1)
    vbox.Add(hbox2)
    window.Add(vbox)
	window.ShowAll()
    
    go func() {
		var d []byte
		wg.Add(1)
		for {
			select {
			case <- quit:
				closeApp()
				wg.Done()
				return
			default:
				mainApp()
			}		
		}
	}()
	
	gtk.Main()
}

func closeApp() {
    
}

func mainApp() {
    
}