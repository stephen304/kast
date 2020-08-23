package backdrop

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"github.com/stephen304/kast/internal"
	"log"
	"os/exec"
	"sync"
)

type Backdrop struct {
	m              sync.Mutex
	cmd            *exec.Cmd
	allocCtxCancel context.CancelFunc
	taskCtx        context.Context
	taskCtxCancel  context.CancelFunc
}

func New(r *gin.RouterGroup, display *internal.DisplayMutex) *Backdrop {
	module := &Backdrop{}

	// TODO It might be better to do this somewhere else
	go display.Assign(module)

	r.POST("/start", func(c *gin.Context) {
		display.Assign(module)
		module.Start()
	})

	r.POST("/stop", func(c *gin.Context) {
		module.Stop()
	})

	r.POST("/prev", func(c *gin.Context) {
		module.m.Lock()
		defer module.m.Unlock()
		if module.taskCtx != nil {
			chromedp.Run(module.taskCtx,
				chromedp.Click(`body > div.pagination a.pagination__link.pagination__link--prev`, chromedp.NodeVisible),
			)
		}
	})

	r.POST("/next", func(c *gin.Context) {
		module.m.Lock()
		defer module.m.Unlock()
		if module.taskCtx != nil {
			chromedp.Run(module.taskCtx,
				chromedp.Click(`body > div.pagination a.pagination__link.pagination__link--next`, chromedp.NodeVisible),
			)
		}
	})

	return module
}

func (module *Backdrop) Start() error {
	module.m.Lock()
	defer module.m.Unlock()

	if module.taskCtx != nil {
		return nil
	}

	opts := []chromedp.ExecAllocatorOption{
		// From: chromedp.DefaultExecAllocatorOptions[:]
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		// chromedp.Headless,

		// After Puppeteer's default behavior.
		chromedp.Flag("disable-background-networking", true),
		// chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
		// chromedp.Flag("disable-background-timer-throttling", true),
		// chromedp.Flag("disable-backgrounding-occluded-windows", true),
		// chromedp.Flag("disable-breakpad", true),
		// chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-default-apps", true),
		// chromedp.Flag("disable-dev-shm-usage", true),
		// chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-features", "site-per-process,TranslateUI,BlinkGenPropertyTrees"),
		// chromedp.Flag("disable-hang-monitor", true),
		// chromedp.Flag("disable-ipc-flooding-protection", true),
		// chromedp.Flag("disable-popup-blocking", true),
		// chromedp.Flag("disable-prompt-on-repost", true),
		// chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-sync", true),
		// chromedp.Flag("force-color-profile", "srgb"),
		// chromedp.Flag("metrics-recording-only", true),
		// chromedp.Flag("safebrowsing-disable-auto-update", true),
		// chromedp.Flag("enable-automation", true), // Causes warning to be shown at the top
		chromedp.Flag("password-store", "basic"),
		chromedp.Flag("use-mock-keychain", true),

		// No UI is best UI
		chromedp.Flag("kiosk", true),
		chromedp.Flag("force-dark-mode", true),
	}

	// Create chrome process
	allocCtx, allocCtxCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	module.allocCtxCancel = allocCtxCancel
	taskCtx, taskCtxCancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	module.taskCtxCancel = taskCtxCancel
	module.taskCtx = taskCtx

	// Start the browser
	err := chromedp.Run(taskCtx,
		chromedp.Navigate(`https://earthview.withgoogle.com/`),
	)
	if err == nil {
		// Using a separate thread doesn't seem to help much here for API responsiveness
		// But at least it seems like it doesn't hurt anything
		go chromedp.Run(taskCtx,
			chromedp.WaitVisible(`body > div.intro a.button.intro__explore`),
			chromedp.Click(`body > div.intro a.button.intro__explore`, chromedp.NodeVisible),
		)
	}

	return err
}

func (module *Backdrop) Stop() error {
	module.m.Lock()
	defer module.m.Unlock()

	if module.taskCtx != nil {
		module.taskCtxCancel()
		module.allocCtxCancel()
		module.taskCtx = nil
	}
	return nil
}
