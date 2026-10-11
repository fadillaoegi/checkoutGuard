// brandassets renders favicon files from the shared vector brand mark.
package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

func main() {
	public := filepath.Join("..", "checkoutguard_FE", "public")
	svg, err := os.ReadFile(filepath.Join(public, "brand", "checkoutguard-mark.svg"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, stop := context.WithTimeout(context.Background(), 25*time.Second)
	defer stop()
	allocator, stopAllocator := chromedp.NewExecAllocator(ctx, chromedp.DefaultExecAllocatorOptions[:]...)
	defer stopAllocator()
	browser, stopBrowser := chromedp.NewContext(allocator)
	defer stopBrowser()
	source, _ := json.Marshal("data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString(svg))
	var rendered map[string]string
	script := fmt.Sprintf(`(async () => { const img=new Image(); img.src=%s; await img.decode(); const result={}; for(const size of [16,32,48,180]) { const c=document.createElement('canvas'); c.width=c.height=size; c.getContext('2d').drawImage(img,0,0,size,size); result[size]=c.toDataURL('image/png').split(',')[1]; } return result; })()`, source)
	if err := chromedp.Run(browser, chromedp.Navigate("about:blank"), chromedp.Evaluate(script, &rendered, func(p *runtime.EvaluateParams) *runtime.EvaluateParams { return p.WithAwaitPromise(true) })); err != nil {
		log.Fatal(err)
	}
	icons := map[int][]byte{}
	for _, size := range []int{16, 32, 48, 180} {
		data, err := base64.StdEncoding.DecodeString(rendered[fmt.Sprint(size)])
		if err != nil {
			log.Fatal(err)
		}
		config, err := png.DecodeConfig(bytes.NewReader(data))
		if err != nil || config.Width != size || config.Height != size {
			log.Fatal("invalid icon raster dimensions")
		}
		name := fmt.Sprintf("favicon-%d.png", size)
		if size == 180 {
			name = "apple-touch-icon.png"
		}
		if err := os.WriteFile(filepath.Join(public, name), data, 0644); err != nil {
			log.Fatal(err)
		}
		icons[size] = data
	}
	// ICO directory with three PNG representations for current browser fallbacks.
	var ico bytes.Buffer
	binary.Write(&ico, binary.LittleEndian, uint16(0))
	binary.Write(&ico, binary.LittleEndian, uint16(1))
	binary.Write(&ico, binary.LittleEndian, uint16(3))
	offset := uint32(6 + 3*16)
	for _, size := range []int{16, 32, 48} {
		ico.Write([]byte{byte(size), byte(size), 0, 0})
		binary.Write(&ico, binary.LittleEndian, uint16(1))
		binary.Write(&ico, binary.LittleEndian, uint16(32))
		binary.Write(&ico, binary.LittleEndian, uint32(len(icons[size])))
		binary.Write(&ico, binary.LittleEndian, offset)
		offset += uint32(len(icons[size]))
	}
	for _, size := range []int{16, 32, 48} {
		ico.Write(icons[size])
	}
	if err := os.WriteFile(filepath.Join(public, "favicon.ico"), ico.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Brand icons generated: 16/32/48px PNG, multi-size ICO and 180px Apple touch icon")
}
