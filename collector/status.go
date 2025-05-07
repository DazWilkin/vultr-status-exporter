package collector

import (
	"log/slog"
	"sync"

	"github.com/DazWilkin/go-probe/probe"
	"github.com/DazWilkin/vultr-status-exporter/api"
	"github.com/prometheus/client_golang/prometheus"
)

// StatusCollector collects metrics about the Vultr service's status API
type StatusCollector struct {
	client *api.Client
	ch     chan<- probe.Status
	logger *slog.Logger

	ServiceAlert   *prometheus.Desc
	Infrastructure *prometheus.Desc
}

// NewStatusCollector returns a new StatusCollector
func NewStatusCollector(ch chan<- probe.Status, l *slog.Logger) *StatusCollector {
	subsystem := "status"
	logger := l.With("collector", subsystem)

	client := api.NewClient()

	return &StatusCollector{
		client: client,
		ch:     ch,
		logger: logger,

		ServiceAlert: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "service_alert"),
			"Vultr Service Alerts",
			[]string{
				"id",
				"region",
				"status",
			},
			nil,
		),
		Infrastructure: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "infrastructure"),
			"Vultr Infrastructure status",
			[]string{
				"region",
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
	logger := c.logger.With("method", "collect")

	// Account for possibility of duplicate ServiceAlert(IDs)
	// See Issue 18
	// https://github.com/DazWilkin/vultr-status-exporter/issues/18
	ids := map[string]bool{}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		alerts, err := c.client.Alerts()
		if err != nil {
			msg := "unable to get Alerts"
			logger.Info(msg,
				"error", err.Error(),
			)
			// Send probe unhealthy probeStatus
			probeStatus := probe.Status{
				Healthy: false,
				Message: msg,
			}
			c.ch <- probeStatus
			return
		}

		// Send probe healthy probeStatus
		probeStatus := probe.Status{
			Healthy: true,
			Message: "ok",
		}
		c.ch <- probeStatus

		for _, serviceAlert := range alerts.ServiceAlerts {
			// Skip ServiceAlert ID if it's been seen
			// ServiceAlerts can be repeated within a region
			// Since these can't be disambiguated, only record them once
			if ids[serviceAlert.ID] {
				logger.Info("StatusCollector:Collect duplicate ServiceAlert",
					"id", serviceAlert.ID,
				)
				continue
			}

			// Record the ServiceAlert ID
			logger.Info("StatusCollector:Collect recording ServiceAlert",
				"id", serviceAlert.ID,
			)
			ids[serviceAlert.ID] = true

			ch <- prometheus.MustNewConstMetric(
				c.ServiceAlert,
				prometheus.CounterValue,
				1.0,
				[]string{
					serviceAlert.ID,
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
			msg := "unable to get Status"
			logger.Info(msg,
				"error", err.Error(),
			)
			// Send probe unhealthy status
			probeStatus := probe.Status{
				Healthy: false,
				Message: msg,
			}
			c.ch <- probeStatus
			return
		}

		// Send probe healthy status
		probeStatus := probe.Status{
			Healthy: true,
			Message: "ok",
		}
		c.ch <- probeStatus

		for ID, region := range status.Regions {
			ch <- prometheus.MustNewConstMetric(
				c.Infrastructure,
				prometheus.GaugeValue,
				float64(len(region.Alerts)),
				[]string{
					ID,
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
	ch <- c.Infrastructure
}
