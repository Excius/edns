package metrics

import "github.com/prometheus/client_golang/prometheus"

type ConnectionMetrics struct {
	ActiveConnections prometheus.Gauge
	ConnectionsOpened prometheus.Counter
	ConnectionsClosed prometheus.Counter
}

type SubscriberMetrics struct {
	MessagesReceived      prometheus.Counter
	MessageHandlingErrors prometheus.Counter
}

type DeliveryMetrics struct {
	MessagesDelivered prometheus.Counter
	DeliveryErrors    prometheus.Counter
	DeliveryDuration  prometheus.Histogram
}

type Metrics struct {
	Connection ConnectionMetrics
	Subscriber SubscriberMetrics
	Delivery   DeliveryMetrics
}

func NewMetrics() *Metrics {
	m := &Metrics{
		Connection: ConnectionMetrics{
			ActiveConnections: prometheus.NewGauge(
				prometheus.GaugeOpts{
					Name: "ws_active_connections",
					Help: "Current active websocket connections.",
				},
			),
			ConnectionsOpened: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "ws_connections_opened_total",
					Help: "Total websocket connections opened.",
				},
			),
			ConnectionsClosed: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "ws_connections_closed_total",
					Help: "Total websocket connections closed.",
				},
			),
		},

		Subscriber: SubscriberMetrics{
			MessagesReceived: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "ws_pubsub_messages_received_total",
					Help: "Total messages received from Redis Pub/Sub.",
				},
			),
			MessageHandlingErrors: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "ws_message_handling_errors_total",
					Help: "Total websocket message handling errors.",
				},
			),
		},

		Delivery: DeliveryMetrics{
			MessagesDelivered: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "ws_messages_delivered_total",
					Help: "Total websocket messages successfully delivered.",
				},
			),
			DeliveryErrors: prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: "ws_delivery_errors_total",
					Help: "Total websocket delivery errors.",
				},
			),
			DeliveryDuration: prometheus.NewHistogram(
				prometheus.HistogramOpts{
					Name:    "ws_delivery_duration_seconds",
					Help:    "Websocket delivery duration.",
					Buckets: prometheus.DefBuckets,
				},
			),
		},
	}

	prometheus.MustRegister(
		m.Connection.ActiveConnections,
		m.Connection.ConnectionsOpened,
		m.Connection.ConnectionsClosed,

		m.Subscriber.MessagesReceived,
		m.Subscriber.MessageHandlingErrors,

		m.Delivery.MessagesDelivered,
		m.Delivery.DeliveryErrors,
		m.Delivery.DeliveryDuration,
	)

	return m
}
