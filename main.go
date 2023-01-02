package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
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
	Display_fb_2_devfb()
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
