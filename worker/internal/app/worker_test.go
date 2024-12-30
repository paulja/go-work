package app_test

import (
	"testing"
	"time"

	"github.com/paulja/go-work/worker/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {
	t.Run("can start and stop worker", func(t *testing.T) {
		w := new(app.Worker)
		var err error
		go func() {
			err = w.Start("1", "testing")
		}()
		time.Sleep(200 * time.Millisecond)
		assert.NoError(t, err, "should be able to start worker")
		assert.NoError(t, w.Stop("1"), "should be able to stop worker")

		// NOTE coverage reports only 50% because it doesn't "see" the
		//      call to the Start method.
	})
}
