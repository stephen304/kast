package media

import (
	"github.com/gin-gonic/gin"
	"github.com/stephen304/kast/internal"
	"log"
	"os"
	"os/exec"
	"sync"
)

type Media struct {
	queue  *queue
	c1     *exec.Cmd
	c2     *exec.Cmd
	m      sync.Mutex
	worker sync.Mutex
}

func New(r *gin.RouterGroup, display *internal.DisplayMutex) *Media {
	module := &Media{
		queue: newQueue(),
	}

	r.GET("/status", func(c *gin.Context) {
		//
	})

	// Adds the url to the end of the queue
	// Also activates the module and starts the process if not active
	r.POST("/enqueue", func(c *gin.Context) {
		url := c.PostForm("url")
		title := url

		// Try to get title
		titleBytes, err := exec.Command("youtube-dl", "--skip-download", "--get-title", "--no-warnings", url).Output()
		if (err == nil) {
			title = string(titleBytes)
		}

		// Add media to queue
		module.queue.enqueue(url, title)

		// Assign module, if it was already assigned, run start to kick off a vlc
		//   thread if vlc is not running
		if !display.Assign(module) {
			module.Start()
		}
	})

	// Immediately plays the selected media followed by nextItems in the queue
	// Activates the module and starts the process if not active
	r.POST("/play", func(c *gin.Context) {
		// url := c.PostForm("url")
		// module.queue.preempt(url)
		// display.Assign(module)
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
		module.Kill()
		module.queue.Next()
		module.Start()
	})

	return module
}

func (module *Media) mediaLoop() {
	// Worker lock lets kill process flush out the worker
	module.worker.Lock()

	for len(module.queue.GetUrl()) > 0 {
		url, title := module.queue.Get()
		log.Printf("[Media] Playing: %s", url)
		module.c1 = exec.Command("youtube-dl", "-o", "-", url)
		module.c2 = exec.Command("cvlc", "-", "vlc://quit", "-f", "--meta-title", title)
		module.c2.Stdin, _ = module.c1.StdoutPipe()
		module.c2.Stdout = os.Stdout // What's this for
		// Run both threads concurrently
		_ = module.c2.Start()
		_ = module.c1.Start()
		// Wait for VLC to finish or be killed
		_ = module.c2.Wait()

		// Process exited
		// Lock while checking process values
		module.m.Lock()
		dead := module.c1 == nil || module.c2 == nil
		module.m.Unlock()
		if dead {
			// Processes were killed, don't touch the queue and just exit
			module.worker.Unlock()
			return
		}
		module.queue.Next()
	}
	// Worker must be unlocked before kill because kill always
	//   waits on the worker lock to ensure thread safety for callers of Kill
	module.worker.Unlock()
	// Kill the process but don't empty the queue
	module.Kill()
}

func (module *Media) Start() error {
	module.m.Lock()
	defer module.m.Unlock()

	if module.c1 == nil && module.c2 == nil {
		go module.mediaLoop()
	}

	return nil
}

func (module *Media) Stop() error {
	module.queue.Empty()
	module.Kill()
	return nil
}

func (module *Media) Kill() {
	module.m.Lock()

	if module.c1 != nil && module.c1.Process != nil {
		module.c1.Process.Kill()
	}
	module.c1 = nil // To ensure it will be restarted later
	if module.c2 != nil && module.c2.Process != nil {
		module.c2.Process.Kill()
	}
	module.c2 = nil

	module.m.Unlock()
	module.worker.Lock()
	module.worker.Unlock()
}
