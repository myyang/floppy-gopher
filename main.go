package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error while run(). Error: %v", err)
		os.Exit(2)
	}

	return
}

func run() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("Can't Init sdl. Error: %v", err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("Fail to Init ttf. Error: %v", err)
	}
	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("Fail to create window. Error: %v", err)
	}

	defer w.Destroy()

	if err := drawTitle(r); err != nil {
		return fmt.Errorf("Can't draw title. Error: %v", err)
	}

	time.Sleep(time.Second * 3)

	scene, err := newScene(r)
	if err != nil {
		return fmt.Errorf("Can't create scene. Error: %v", err)
	}
	defer scene.destroy()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	select {
	case err := <-scene.run(ctx, r):
		return err
	case <-time.After(time.Second * 4):
		return nil
	}

	return nil
}

func drawTitle(r *sdl.Renderer) error {
	r.Clear()

	font, err := ttf.OpenFont("fonts/floppy.ttf", 20)
	if err != nil {
		return fmt.Errorf("Can't open 'fonts/floppy.ttf'. Error: %v", err)
	}
	defer font.Close()

	surface, err := font.RenderUTF8_Solid("Flappy Gopher", sdl.Color{R: 255, G: 100, B: 0, A: 255})
	if err != nil {
		return fmt.Errorf("Can't render title. Error: %v", err)
	}
	defer surface.Free()

	texture, err := r.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("Can't create texture. Error: %v", err)
	}
	defer texture.Destroy()

	if err := r.Copy(texture, nil, nil); err != nil {
		return fmt.Errorf("Could not copy texture to render. Error: %v", err)
	}

	r.Present()

	return nil
}

func drawBackground(r *sdl.Renderer) error {
	r.Clear()

	texture, err := img.LoadTexture(r, "img/background.png")
	if err != nil {
		return fmt.Errorf("Error to load background. Error: %v", err)
	}
	defer texture.Destroy()

	if err := r.Copy(texture, nil, nil); err != nil {
		return fmt.Errorf("Could not copy bird to render. Error: %v", err)
	}

	r.Present()

	return nil
}
