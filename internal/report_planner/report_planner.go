package report_planner

import (
	"container/list"
	"crypto/rand"
	"golang.org/x/net/context"
	"log/slog"
	"math/big"
	"time"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const traceIdLength = 16
const queueSleepTime = 10 * time.Second
const reportGeneratorChannelBuffer = 100

type AverageLoadingStorage interface {
	GetAverageLoadingTime(
		// duration for 1 day in load date range
		ctx context.Context,
		reportName string,
	) (
		time.Duration,
		error,
	)
}

type ReportPlanner struct {
	reportsLocalQueue     *list.List
	averageLoadingStorage AverageLoadingStorage
	log                   *slog.Logger
}

func NewReportPlanner(
	logger *slog.Logger,
	averageLoadingStorage AverageLoadingStorage,
) *ReportPlanner {
	queue := list.New()
	return &ReportPlanner{
		reportsLocalQueue:     queue,
		averageLoadingStorage: averageLoadingStorage,
		log:                   logger,
	}
}

type reportQueueItem struct {
	reportName      string
	reportFunction  func() error
	estimatedDate   time.Time
	loadingDuration time.Duration
	traceId         string
}

func (i *reportQueueItem) GetReserveTime() time.Duration {
	now := time.Now().UTC()
	return i.estimatedDate.Sub(now.Add(i.loadingDuration))
}

func (p *ReportPlanner) Add(
	reportName string,
	reportFunc func() error,
	estimatedDate time.Time,
	dateFrom time.Time,
	dateTo time.Time,
) {
	traceId := p.generateTraceId()
	p.log.Info(
		"add report to queue",
		"reportName",
		reportName,
		"DateFrom",
		dateFrom,
		"DateTo",
		dateTo,
		"estimatedDate",
		estimatedDate,
		"traceId",
		traceId,
	)

	ctx := context.Background()
	loadingStatDuration, err := p.averageLoadingStorage.GetAverageLoadingTime(
		ctx,
		reportName,
	)
	if err != nil {
		return
	}

	loadingDuration := time.Duration(dateTo.Add(24*time.Hour).Sub(dateFrom).Hours()/24) * loadingStatDuration

	item := &reportQueueItem{
		reportName:      reportName,
		reportFunction:  reportFunc,
		estimatedDate:   estimatedDate.UTC(),
		loadingDuration: loadingDuration,
		traceId:         traceId,
	}
	p.addReportItemToQueue(item)
}

func (p *ReportPlanner) Get() chan *reportQueueItem {
	ch := make(
		chan *reportQueueItem,
		reportGeneratorChannelBuffer,
	)

	go func(log *slog.Logger) {
		for {
			currentItem := p.reportsLocalQueue.Front()
			if currentItem == nil {
				log.Info("queue is empty, sleep")
				time.Sleep(queueSleepTime)
				continue
			}
			ch <- p.reportsLocalQueue.Remove(currentItem).(*reportQueueItem)
		}
	}(p.log)
	return ch
}

func (p *ReportPlanner) addReportItemToQueue(item *reportQueueItem) {
	rep := p.reportsLocalQueue.Front()

	if rep == nil {
		p.reportsLocalQueue.PushBack(item)
		p.log.Info(
			"Report Added to queue",
			"TraceId",
			item.traceId,
		)
		return
	}

	for i := 0; i <= p.reportsLocalQueue.Len(); i++ {
		currRep := rep.Value.(*reportQueueItem)
		if currRep.GetReserveTime() <= item.GetReserveTime() {
			next := rep.Next()
			if next == nil {
				p.reportsLocalQueue.InsertAfter(
					item,
					rep,
				)
				p.log.Info(
					"Report Added to queue",
					"TraceId",
					item.traceId,
				)
				break
			} else {
				rep = next
			}
		} else {
			p.reportsLocalQueue.InsertBefore(
				item,
				rep,
			)
			p.log.Info(
				"Report Added to queue",
				"TraceId",
				item.traceId,
			)
			break
		}
	}
}

func (p *ReportPlanner) generateTraceId() string {
	b := make(
		[]byte,
		traceIdLength,
	)
	for i := range b {
		num, _ := rand.Int(
			rand.Reader,
			big.NewInt(int64(len(charset))),
		)
		b[i] = charset[num.Int64()]
	}
	return string(b)
}

func (p *ReportPlanner) GetAllSequence() []*reportQueueItem {
	var sequence []*reportQueueItem
	currentItem := p.reportsLocalQueue.Front()
	for {
		if currentItem == nil {
			return sequence
		}
		item := currentItem.Value.(*reportQueueItem)
		sequence = append(
			sequence,
			item,
		)
		currentItem = currentItem.Next()
	}
}
