package metrics

import "github.com/prometheus/client_golang/prometheus"

type ConsumerMetrics struct {
	MessagesReceived        prometheus.Counter
	MessagesAcknowledged    prometheus.Counter
	MessageProcessingErrors prometheus.Counter
}

type ProcessorMetrics struct {
	NotificationsProcessed prometheus.Counter
	NotificationsCompleted prometheus.Counter
	NotificationsFailed    prometheus.Counter

	DeliveriesProcessed prometheus.Counter
	DeliveriesCompleted prometheus.Counter
	DeliveriesFailed    prometheus.Counter

	DeliveryRetries prometheus.Counter
	DLQMessages     prometheus.Counter

	NotificationProcessingDuration prometheus.Histogram
	DeliveryDuration               prometheus.Histogram
}

type RecoveryMetrics struct {
	RecoveryScans     prometheus.Counter
	RecoveredMessages prometheus.Counter
	RecoveryErrors    prometheus.Counter

	RecoveryProcessingDuration prometheus.Histogram
}

type Metrics struct {
	Consumer  ConsumerMetrics
	Processor ProcessorMetrics
	Recovery  RecoveryMetrics
}

func NewMetrics() *Metrics {
	m := &Metrics{
		Consumer: ConsumerMetrics{
			MessagesReceived: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "worker_stream_messages_received_total",
					Help: "Total Redis Stream messages received by the worker.",
				},
			),
			MessagesAcknowledged: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "worker_stream_messages_acknowledged_total",
					Help: "Total Redis Stream messages acknowledged by the worker.",
				},
			),
			MessageProcessingErrors: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "worker_stream_processing_errors_total",
					Help: "Total message processing errors.",
				},
			),
		},

		Processor: ProcessorMetrics{
			NotificationsProcessed: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "worker_notifications_processed_total",
					Help: "Total notifications processed.",
				},
			),
			NotificationsCompleted: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "worker_notifications_completed_total",
					Help: "Total notifications completed.",
				},
			),
			NotificationsFailed: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "worker_notifications_failed_total",
					Help: "Total notifications failed.",
				},
			),
			DeliveriesProcessed: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "worker_deliveries_processed_total",
					Help: "Total deliveries processed.",
				},
			),
			DeliveriesCompleted: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "worker_deliveries_completed_total",
					Help: "Total deliveries completed.",
				},
			),
			DeliveriesFailed: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "worker_deliveries_failed_total",
					Help: "Total permanently failed deliveries.",
				},
			),
			DeliveryRetries: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "worker_delivery_retries_total",
					Help: "Total delivery retries.",
				},
			),
			DLQMessages: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "worker_dlq_messages_total",
					Help: "Total messages published to the DLQ.",
				},
			),
			NotificationProcessingDuration: prometheus.NewHistogram(
				prometheus.HistogramOpts{
					Name:    "worker_notification_processing_duration_seconds",
					Help:    "Notification processing duration.",
					Buckets: prometheus.DefBuckets,
				},
			),
			DeliveryDuration: prometheus.NewHistogram(
				prometheus.HistogramOpts{
					Name:    "worker_delivery_duration_seconds",
					Help:    "Delivery duration.",
					Buckets: prometheus.DefBuckets,
				},
			),
		},

		Recovery: RecoveryMetrics{
			RecoveryScans: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "worker_recovery_scans_total",
					Help: "Total recovery scans executed.",
				},
			),
			RecoveredMessages: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "worker_recovered_messages_total",
					Help: "Total messages recovered from the Pending Entries List.",
				},
			),
			RecoveryErrors: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "worker_recovery_errors_total",
					Help: "Total recovery errors.",
				},
			),
			RecoveryProcessingDuration: prometheus.NewHistogram(
				prometheus.HistogramOpts{
					Name:    "worker_recovery_processing_duration_seconds",
					Help:    "Time spent processing recovered messages.",
					Buckets: prometheus.DefBuckets,
				},
			),
		},
	}

	prometheus.MustRegister(

		// Consumer
		m.Consumer.MessagesReceived,
		m.Consumer.MessagesAcknowledged,
		m.Consumer.MessageProcessingErrors,

		// Processor
		m.Processor.NotificationsProcessed,
		m.Processor.NotificationsCompleted,
		m.Processor.NotificationsFailed,

		m.Processor.DeliveriesProcessed,
		m.Processor.DeliveriesCompleted,
		m.Processor.DeliveriesFailed,

		m.Processor.DeliveryRetries,
		m.Processor.DLQMessages,

		m.Processor.NotificationProcessingDuration,
		m.Processor.DeliveryDuration,

		// Recovery
		m.Recovery.RecoveryScans,
		m.Recovery.RecoveredMessages,
		m.Recovery.RecoveryErrors,
		m.Recovery.RecoveryProcessingDuration,
	)

	return m
}
