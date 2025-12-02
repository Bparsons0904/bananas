package seeder

import (
	"fmt"
	"time"

	"github.com/schollz/progressbar/v3"
)

// ProgressTracker manages progress bars for seeding operations
type ProgressTracker struct {
	currentBar *progressbar.ProgressBar
	startTime  time.Time
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker() *ProgressTracker {
	return &ProgressTracker{
		startTime: time.Now(),
	}
}

// StartTable starts tracking progress for a table
func (pt *ProgressTracker) StartTable(tableName string, total int) {
	pt.currentBar = progressbar.NewOptions(total,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(40),
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan]%s[reset]", tableName)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("records"),
	)
}

// Add increments the progress bar
func (pt *ProgressTracker) Add(count int) error {
	if pt.currentBar != nil {
		return pt.currentBar.Add(count)
	}
	return nil
}

// Finish completes the current progress bar
func (pt *ProgressTracker) Finish() error {
	if pt.currentBar != nil {
		return pt.currentBar.Finish()
	}
	return nil
}

// Elapsed returns the time since tracking started
func (pt *ProgressTracker) Elapsed() time.Duration {
	return time.Since(pt.startTime)
}

// FormatDuration formats a duration for display
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	return fmt.Sprintf("%.1fh", d.Hours())
}
