package backdrop

import (
	"github.com/gin-gonic/gin"
	"github.com/stephen304/kast/internal"
	// "math/rand"
	"os/exec"
	"regexp"
	// "strings"
	"sync"
	// "time"
	"fmt"
	"log"
)

const NAME = "BACKDROP"

type Backdrop struct {
	m   sync.Mutex
	cmd *exec.Cmd
}

func New(r *gin.RouterGroup, display *internal.DisplayMutex) *Backdrop {
	module := &Backdrop{}

	display.Assign(module)

	r.POST("/start", func(c *gin.Context) {
		// url := c.PostForm("url")
		// module.queue.enqueue(url)
		display.Assign(module)
		module.m.Lock()
		if module.cmd == nil {
			module.Start()
		}
		module.m.Unlock()
	})

	r.POST("/stop", func(c *gin.Context) {
		// url := c.PostForm("url")
		// module.queue.enqueue(url)
		module.m.Lock()
		if module.cmd != nil {
			module.Stop()
		}
		module.m.Unlock()
	})

	return module
}

func (module *Backdrop) GetName() string {
	return NAME
}

func (module *Backdrop) Start() error {
	if module.cmd != nil {
		return nil
	}
	backgroundsList, err := exec.Command("curl", "-s", "-o", "-", "https://raw.githubusercontent.com/dconnolly/chromecast-backgrounds/master/README.md").Output()
	if err != nil {
		module.cmd = nil
		return err
	}
	// backgrounds := strings.Split(string(backgroundsList), "\n")
	backgrounds := string(backgroundsList)
	// rand.Seed(time.Now().UnixNano())
	// rand.Shuffle(len(backgrounds), func(i, j int) { backgrounds[i], backgrounds[j] = backgrounds[j], backgrounds[i] })

	re := regexp.MustCompile(`https?://[^)]+`)
	matches := re.FindAllStringSubmatch(string(backgrounds), -1)

	var cleanMatches []string
	for _, x := range matches {
		cleanMatches = append(cleanMatches, x[0]) // note the = instead of :=
	}
	fmt.Printf("%+v\n", cleanMatches)
	flags := []string{"-F", "-p", "-z"}
	flags = append(flags, cleanMatches...)
	// log.Printf("cleanMatches: ", cleanMatches)
	feh := exec.Command("feh", flags...)
	// err = feh.Start()
	output, err := feh.CombinedOutput()

	log.Printf("Output:\n%s", output)
	if err != nil {
		module.cmd = nil
		return err
	}

	module.cmd = feh

	return nil
}

func (module *Backdrop) Stop() error {
	var err error = nil
	if module.cmd != nil {
		err = module.cmd.Process.Kill() // Immediate
		module.cmd = nil
	}
	return err
}
