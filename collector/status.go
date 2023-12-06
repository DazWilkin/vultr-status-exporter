package collector

import (
	"log/slog"
	"sync"

	"github.com/DazWilkin/vultr-status-exporter/api"
	"github.com/prometheus/client_golang/prometheus"
)

// StatusCollector collects metrics about the Vultr service's status API
type StatusCollector struct {
	client *api.Client

	ServiceAlert *prometheus.Desc
	Region       *prometheus.Desc
}

// NewStatusCollector returns a new StatusCollector
func NewStatusCollector() *StatusCollector {
	subsystem := "status"

	client := api.NewClient()

	return &StatusCollector{
		client: client,

		ServiceAlert: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "service_alert"),
			"Alert",
			[]string{
				"region",
				"status",
			},
			nil,
		),
		Region: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "region"),
			"Region",
			[]string{
				"location",
				"country",
				"country_name",
			},
			nil,
		),
	}
}

// Collect implements Prometheus' Collector interface and is used to collect metrics
func (c *StatusCollector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		alerts, err := c.client.Alerts()
		if err != nil {
			slog.Info("StatusCollector:Collect",
				"error", err.Error(),
			)
			return
		}

		for _, serviceAlert := range alerts.ServiceAlerts {
			ch <- prometheus.MustNewConstMetric(
				c.ServiceAlert,
				prometheus.CounterValue,
				1.0,
				[]string{
					serviceAlert.Region,
					serviceAlert.Status,
				}...,
			)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		status, err := c.client.Status()
		if err != nil {
			slog.Info("StatusCollector:Collect",
				"error", err.Error(),
			)
			return
		}

		for _, region := range status.Regions {
			ch <- prometheus.MustNewConstMetric(
				c.Region,
				prometheus.GaugeValue,
				float64(len(region.Alerts)),
				[]string{
					region.Location,
					region.Country,
					region.CountryName,
				}...,
			)
		}
	}()

	wg.Wait()
}

// Describe implements Prometheus' Collector interface and is used to describe metrics
func (c *StatusCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.ServiceAlert
	ch <- c.Region
}
