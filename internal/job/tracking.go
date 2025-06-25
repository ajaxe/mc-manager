package job

import (
	"sync"
	"time"

	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/labstack/echo/v4"
)

var thresholdMinutes = 10
var warnThreshold = time.Duration(thresholdMinutes) * time.Minute

// observer is a struct that manages the current play timer item and its associated timer.
type observer struct {
	logger  echo.Logger
	mu      sync.Mutex
	current *models.PlayTimerItem
	timer   *time.Timer
	endDt   time.Time
	C       <-chan time.Time
}

// setCurrent assigns the new play timer to the observer.
func (j *observer) setCurrent(item *models.PlayTimerItem) (err error) {
	dt, err := time.Parse(time.RFC3339, item.EndDate)
	if err != nil {
		return
	}

	j.mu.Lock()
	defer j.mu.Unlock()
	j.current = item
	if j.timer != nil {
		j.timer.Stop()
		j.timer = nil
		j.C = nil
	}

	j.endDt = dt

	d := nextTick(j.endDt)
	if d <= 0 {
		j.logger.Infof("cannot set current item: play time exceeded Now()")
		expirePlayTimer(item, j.logger)
	}
	j.logger.Infof("setting first tick after %v minutes", int(d.Minutes()))
	j.timer = time.NewTimer(d)
	j.C = j.timer.C

	return
}

// setNextTick sets the timer for the next tick based on the end date of the current play timer item.
func (j *observer) setNextTick() (done bool) {
	j.mu.Lock()
	defer j.mu.Unlock()
	d := nextTick(j.endDt)
	if d <= 0 {
		j.logger.Infof("play time exceeded Now()")
		expirePlayTimer(j.current, j.logger)
		return true
	}
	j.logger.Infof("setting next tick after %v minutes", int(d.Minutes()))
	j.timer = time.NewTimer(d)
	j.C = j.timer.C
	return false
}

// sendMessage processes the current play timer item and sends a message to players if necessary.
func (j *observer) remaining() time.Duration {
	j.mu.Lock()
	defer j.mu.Unlock()
	diff := j.endDt.Sub(time.Now().UTC())

	j.logger.Infof("remaining duration: %v", diff)

	return diff
}

func (j *observer) stop() {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.timer != nil {
		j.timer.Stop()
		j.timer = nil
		j.C = nil
	}
	j.logger.Info("Stopping observer timer")
}

// nextTick returns the duration until the next tick when an action is to be taken, like sending message to players
func nextTick(t time.Time) (d time.Duration) {
	if time.Now().UTC().After(t) {
		d = 0
	}

	diff := t.Sub(time.Now().UTC())

	if diff > warnThreshold {
		d = diff - warnThreshold
		return
	}

	if diff > time.Minute {
		d = time.Minute
		return
	}

	d = diff

	return
}
