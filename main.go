package main

import (
	"os"
	"fmt"
	"image"
	"image/color"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
    "golang.org/x/image/math/fixed"
	"github.com/eiannone/keyboard"
)

const color_channel  int = 4

var (
	fb_width, fb_width_stride, fb_height int
	fbframe *image.NRGBA;
)


func main() {
	Prep_fb_img()

	for y:=0; y<1000; y++ {
		fbframe.Set(10, y, color.NRGBA{255, 128, 128, 0})
	}


	//test
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press ESC to quit")

	for xcursor:=10; ; {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		putString(xcursor, 100, string(char))
		xcursor += 7
		Display_fb_2_devfb()

		if key == keyboard.KeyCtrlZ {
			break
		}
	}
}


func GetSizeOfScreen()  {
	bytes, err := os.ReadFile("/sys/class/graphics/fb0/virtual_size")
	if err != nil {
		fmt.Printf("Something went wrong: %v\n", err)
		panic(err)
	}
	fmt.Sscanf(string(bytes), "%v,%v\n",  &fb_width, &fb_height)


	bytes, err = os.ReadFile("/sys/class/graphics/fb0/stride")
	if err != nil {
		fmt.Printf("Something went wrong: %v\n", err)
		panic(err)
	}
	fmt.Sscanf(string(bytes), "%v\n",  &fb_width_stride)
	fb_width_stride /= 4
}

func Prep_fb_img(){
	GetSizeOfScreen();
	rect := image.Rect(0, 0, fb_width_stride, fb_height)
	
	fbframe = &image.NRGBA{
		Pix: make([]uint8, rect.Dx() * rect.Dy() * color_channel), 
		Stride: 4 * fb_width_stride,
		Rect: rect,
	}
}

func Display_fb_2_devfb(){
	err := os.WriteFile("/dev/fb0", fbframe.Pix, 0644)
	if err!=nil {
		fmt.Printf("Something wrong with writing to framebuffer dev, Err: %v", err)
		panic(err)
	}
	
}

//func putString(img *image.NRGBA, x, y int, label string) {
func putString(x, y int, label string) {
    col := color.NRGBA{200, 100, 0, 255}
    point := fixed.Point26_6{fixed.I(x), fixed.I(y)}

    d := &font.Drawer{
        Dst:  fbframe,
        Src:  image.NewUniform(col),
        Face: basicfont.Face7x13,
        Dot:  point,
    }
    d.DrawString(label)
}

