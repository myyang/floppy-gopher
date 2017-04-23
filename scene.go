package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

type scene struct {
	bg    *sdl.Texture
	bird  *bird
	pipes *pipes
}

func newScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, "img/background.png")
	if err != nil {
		return nil, fmt.Errorf("Error to load background. Error: %v", err)
	}
	bird, err := newBird(r)
	if err != nil {
		return nil, fmt.Errorf("Fail to create new bird. Error: %v", err)
	}
	pipes, err := newPipes(r)
	if err != nil {
		return nil, err
	}

	return &scene{bg: bg, bird: bird, pipes: pipes}, nil
}

func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case event := <-events:
				if done := s.handleEvent(event); done {
					return
				}
			case <-tick:
				s.update()

				if s.bird.isDead() {
					drawTitle(r, "Game Over")
					time.Sleep(time.Second * 2)
					s.restart()
				}

				if err := s.print(r); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseButtonEvent:
		s.bird.jump()
	default:
		//log.Printf("Event %+v", event)
	}
	return false
}

func (s *scene) update() {
	s.bird.update()
	s.pipes.update()
	s.pipes.touch(s.bird)
}

func (s *scene) restart() {
	s.bird.restart()
	s.pipes.restart()
}

func (s *scene) print(r *sdl.Renderer) error {
	r.Clear()

	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("Could not copy background to render. Error: %v", err)
	}

	if err := s.bird.paint(r); err != nil {
		return err
	}
	if err := s.pipes.paint(r); err != nil {
		return err
	}

	r.Present()

	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.destory()
	s.pipes.destory()
}
