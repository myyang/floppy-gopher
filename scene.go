package main

import (
	"context"
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

type scene struct {
	time  int
	bg    *sdl.Texture
	birds []*sdl.Texture
}

func newScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, "img/background.png")
	if err != nil {
		return nil, fmt.Errorf("Error to load background. Error: %v", err)
	}
	var birds []*sdl.Texture
	for i := 1; i < 5; i++ {
		bird, err := img.LoadTexture(r, fmt.Sprintf("img/frame-%d.png", i))
		if err != nil {
			return nil, fmt.Errorf("Error to load bird. Error: %v", err)
		}
		birds = append(birds, bird)
	}

	return &scene{bg: bg, birds: birds}, nil
}

func (s *scene) run(ctx context.Context, r *sdl.Renderer) <-chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)
		for range time.Tick(10 * time.Millisecond) {
			select {
			case <-ctx.Done():
				return
			default:
				if err := s.print(r); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (s *scene) print(r *sdl.Renderer) error {
	s.time++
	r.Clear()

	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("Could not copy background to render. Error: %v", err)
	}

	birdRect := &sdl.Rect{X: 50, Y: 300 - 43/2, W: 50, H: 43}
	if err := r.Copy(s.birds[s.time/10%len(s.birds)], nil, birdRect); err != nil {
		return fmt.Errorf("Could not copy bird to render. Error: %v", err)
	}
	r.Present()

	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
}
