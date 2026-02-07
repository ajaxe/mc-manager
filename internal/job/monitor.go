package job

import (
	"context"
	"fmt"
	"time"

	"github.com/ajaxe/mc-manager/internal/config"
	"github.com/ajaxe/mc-manager/internal/db"
	"github.com/ajaxe/mc-manager/internal/gameserver"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var currentObserver = &observer{}

var inputCh = make(chan *models.PlayTimerItem, 1)

func QueueJob(item *models.PlayTimerItem) {
	inputCh <- item
}

func StopCurrentPlayTimer() {
	currentObserver.stop()
}

func StartMonitor(ctx context.Context, l echo.Logger, client *db.Client) {
	currentObserver.logger = l
	currentObserver.db = client

	i, err := client.ActivePlayTimer()

	ops := gameserver.NewGameServerOperations(l, config.LoadAppConfig())

	if err != nil {
		l.Error("on startup, failed to get active play timer:", err)
	}

	if i != nil {
		currentObserver.setCurrent(i)
	}

	for {
		select {
		case i := <-inputCh:
			l.Infof("received job item: %v: end date: %v", i.ID, i.EndDate)
			currentObserver.setCurrent(i)
		case <-currentObserver.C:
			sendMessage(ops, currentObserver.remaining())
			if n := currentObserver.setNextTick(); n {
				ops.Message("Shutting down server ...")
				time.Sleep(2 * time.Second)
				ops.StopAll()
				currentObserver.stop()
			}
		case <-ctx.Done():
			currentObserver.stop()
			l.Info("Stopping job monitor")
			return
		}
	}
}

func sendMessage(c *gameserver.GameServerOperations, d time.Duration) {
	u := "minutes"
	v := int(d.Minutes())

	if v <= 0 {
		u = "seconds"
		v = int(d.Seconds())
	}

	m := fmt.Sprintf("Remaining play time: %v %s", v, u)

	_ = c.Message(m)
}

func (j *observer) expirePlayTimer(p *models.PlayTimerItem) {
	if p == nil {
		return
	}

	p.IsActive = false
	p.LastUpdateDate = time.Now().UTC().Format(time.RFC3339)

	id, err := bson.ObjectIDFromHex(p.ID)

	if err != nil {
		j.logger.Errorf("failed to parse play timer ID: %v: %v", p.ID, err)
		return
	}
	if err := j.db.UpdatePlayTimerByID(id, p); err != nil {
		j.logger.Errorf("failed to update play timer ID: %v: %v", p.ID, err)
	} else {
		j.logger.Infof("Play timer ID: %v expired and updated successfully", p.ID)
	}
}
