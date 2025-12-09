package service

import (
	"context"
	"sync"
	"time"

	"toll/api/repository"
	"toll/api/types"

	"toll/internal/database"
	log "toll/internal/log"
)

const bufferSize = 2000
const workesCount = 5
const workerBatchSizeLimit = 2000
const workerTimeoutTrigger = 30 * time.Second

type timeRangeFee struct {
	startHour, startMinute int
	endHour, endMinute     int
	fee                    int
}

var holidays = map[string]struct{}{
	"2025-01-01": {},
	"2025-12-25": {},
}

// Toll fee prices hour list.
// {startHr, startMin, endHr, endMin, price}.
var tollFees = []timeRangeFee{
	{6, 0, 6, 29, 8},
	{6, 30, 6, 59, 13},
	{7, 0, 7, 59, 18},
	{8, 0, 8, 29, 13},
	{8, 30, 14, 59, 8},
	{15, 0, 15, 29, 13},
	{15, 30, 16, 59, 18},
	{17, 0, 17, 59, 13},
	{18, 0, 18, 29, 8},
}

var maxDailyFee = 60

type (
	// BillingService interface with method definitions.
	BillingService interface {
		TriggerFor(license string)
	}

	billing struct {
		log log.Logger

		cancel      context.CancelFunc
		wg          sync.WaitGroup
		licenseCh   chan string
		workerCount int

		events repository.TollEventRepository
	}
)

// BillingWorkers func returns new BillingService with workers started.
func BillingWorkers(workerCount int, bufferSize int) BillingService {
	db := database.Get()
	ctx, cancel := context.WithCancel(context.Background())

	svc := &billing{
		log: log.WithField(types.LogComponent, "billing"),

		events:      repository.TollEvent(db),
		licenseCh:   make(chan string, bufferSize),
		workerCount: workerCount,
		cancel:      cancel,
	}

	for i := 0; i < workerCount; i++ {
		svc.wg.Add(1)

		go svc.worker(ctx, i+1, workerBatchSizeLimit, workerTimeoutTrigger)
	}

	return svc
}

// TriggerFor sends a license to the processing channel.
func (svc *billing) TriggerFor(license string) {
	svc.licenseCh <- license
}

func (svc *billing) billLicenses(ctx context.Context, licenses []string) error {
	if len(licenses) == 0 {
		return nil
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// fetch ALL events for the day (all licenses)
	events, err := svc.events.GetAll(ctx, startOfDay, licenses)
	if err != nil {
		return err
	}

	// group events by license plate
	eventsByLicense := make(map[string][]*types.TollEvent)
	for _, ev := range events {
		eventsByLicense[ev.LicensePlate] = append(eventsByLicense[ev.LicensePlate], ev)
	}

	for _, license := range licenses {
		levents := eventsByLicense[license]

		totalFee := svc.calculateDailyFee(levents)

		// write fee row for this license
		err := svc.events.UpdateDailyFee(
			ctx,
			types.DailyFee{
				Date:         startOfDay,
				LicensePlate: license,
				Fee:          totalFee,
			},
		)
		if err != nil {
			svc.log.Errorf("Billing for %s failed: %v", license, err)
			return err
		}

		svc.log.Infof(
			"Billing %s: %d events, total fee = %d SEK (took %s)",
			license, len(levents), totalFee, time.Since(now),
		)
	}

	return nil
}

func IsTollFreeDate(t time.Time) bool {
	weekday := t.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}

	dateStr := t.Format("2006-01-02")
	if _, exists := holidays[dateStr]; exists {
		return true
	}

	return false
}

func priceForEvent(event *types.TollEvent) int {
	if IsTollFreeDate(event.EventStart) || event.IsTollFree() {
		return 0
	}

	hour := event.EventStart.Hour()
	minute := event.EventStart.Minute()

	for _, tf := range tollFees {
		if (hour > tf.startHour || (hour == tf.startHour && minute >= tf.startMinute)) &&
			(hour < tf.endHour || (hour == tf.endHour && minute <= tf.endMinute)) {
			return tf.fee
		}
	}

	return 0
}

func (svc *billing) calculateDailyFee(events []*types.TollEvent) int {
	if len(events) == 0 {
		svc.log.Debug("no events; total fee = 0")

		return 0
	}

	windowStart := events[0].EventStart
	maxFeeInWindow := priceForEvent(events[0])

	svc.log.Debugf("start window at %s with event start at %s fee=%d", windowStart, events[0].EventStart, maxFeeInWindow)

	total := 0

	for _, event := range events[1:] {
		fee := priceForEvent(event)
		diff := event.EventStart.Sub(windowStart)

		if diff < time.Hour {
			// Inside the same 1hr window.
			if fee > maxFeeInWindow {
				svc.log.Debugf("event with start at %s fee=%d replaces previous max fee=%d in same window", event.EventStart, fee, maxFeeInWindow)

				maxFeeInWindow = fee
			} else {
				svc.log.Debugf("event with start at %s fee=%d ignored windowStart=%s, maxFeeInWindow=%d", event.EventStart, fee, windowStart, maxFeeInWindow)
			}
		} else {
			// Window ended add the max fee and start a new one.
			svc.log.Debugf("closing window starting at %s: adding fee=%d", windowStart, maxFeeInWindow)

			total += maxFeeInWindow

			svc.log.Debugf("new window start at %s with event start at %s fee=%d", windowStart, event.EventStart, fee)

			windowStart = event.EventStart
			maxFeeInWindow = fee
		}
	}

	svc.log.Debugf("closing final window starting at %s: adding fee=%d", windowStart, maxFeeInWindow)

	total += maxFeeInWindow

	// Apply daily maximum cap.
	if total > maxDailyFee {
		svc.log.Debugf("total fee %d exceeds maxDailyFee=%d ; total fee = daily maximum cap", total, maxDailyFee)

		return maxDailyFee
	}

	svc.log.Debugf("total fee for day = %d", total)

	return total
}

// worker listens to the license channel and triggers billing in batches.
func (svc *billing) worker(ctx context.Context, workerID int, batchSize int, flushInterval time.Duration) {
	defer svc.wg.Done()

	batch := make([]string, 0, batchSize)

	flush := func() {
		if len(batch) == 0 {
			return
		}

		svc.log.Infof("Worker %d: flushing %d licenses", workerID, len(batch))

		if err := svc.billLicenses(ctx, batch); err != nil {
			svc.log.Infof("Worker %d: error processing license %s: %v", workerID, batch, err)
		}

		batch = batch[:0]
	}

	timer := time.NewTimer(flushInterval)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			svc.log.Infof("Worker %d stopping due to context cancellation", workerID)

			flush()

			return

		case license, ok := <-svc.licenseCh:
			if !ok {
				svc.log.Infof("Worker %d exiting, channel closed", workerID)

				flush()

				return
			}

			batch = append(batch, license)

			if len(batch) >= batchSize {
				flush()

				if !timer.Stop() {
					<-timer.C
				}

				timer.Reset(flushInterval)
			}

		case <-timer.C:
			// flush because of timeout
			flush()

			timer.Reset(flushInterval)
		}
	}
}
