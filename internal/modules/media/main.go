package media

import (
	"github.com/gin-gonic/gin"
	"github.com/stephen304/kast/internal"
	"log"
	"os"
	"os/exec"
	"sync"
)

const NAME = "MEDIA"

type Media struct {
	queue   *queue
	c1      *exec.Cmd
	c2      *exec.Cmd
	m       sync.Mutex
	worker       sync.Mutex
	running bool
}

func New(r *gin.RouterGroup, display *internal.DisplayMutex) *Media {
	module := &Media{
		queue:   newQueue(),
		running: false,
	}

	r.POST("/load", func(c *gin.Context) {
		url := c.PostForm("url")
		module.queue.enqueue(url)
		display.Assign(module)
	})

	// r.POST("/pause", func(c *gin.Context) {
	//
	// })
	//
	// r.POST("/unpause", func(c *gin.Context) {
	//
	// })
	//
	r.POST("/prev", func(c *gin.Context) {
		module.Kill()
		module.queue.Prev()
		module.Start()
	})

	r.POST("/next", func(c *gin.Context) {

	})

	return module
}

func (module *Media) GetName() string {
	return NAME
}

func (module *Media) mediaLoop() {
	module.worker.Lock()
	defer module.worker.Unlock()

	for len(module.queue.Get()) > 0 {
		url := module.queue.Get()
		log.Printf("[Media] playing: %s", url)
		module.c1 = exec.Command("youtube-dl", "-o", "-", url)
		module.c2 = exec.Command("cvlc", "-", "vlc://quit", "-f", "--no-video-title-show")
		module.c2.Stdin, _ = module.c1.StdoutPipe()
		module.c2.Stdout = os.Stdout // What's this for
		_ = module.c2.Start()
		_ = module.c1.Run()
		_ = module.c2.Wait()

		// Process exited
		// Lock is to keep this worker from running until any kill actions are complete
		module.m.Lock()
		if module.running == false {
			// Processes were killed
			module.m.Unlock()
			return
		}
		module.queue.Next()
		module.m.Unlock()
	}
	module.Stop()
}

func (module *Media) Start() error {
	module.running = true
	go module.mediaLoop()
	return nil
}

func (module *Media) Stop() error {
	log.Println("[Media] Stopping...")
	module.queue.Empty()
	module.Kill()
	return nil
}

func (module *Media) Kill() {
	module.m.Lock()

	if module.c1 != nil && module.c1.Process != nil {
		module.c1.Process.Kill()
	}
	if module.c2 != nil && module.c2.Process != nil {
		module.c2.Process.Kill()
	}
	module.running = false
	module.m.Unlock()
	// Ensure all workers have exited
	module.worker.Lock()
	module.worker.Unlock()
}
