package kswitch

import (
	"time"
)

type Event struct {
	Version           string  `json:"version"`
	GroupKey          string  `json:"groupKey"`
	TruncatedAlerts   int     `json:"truncatedAlerts"`
	Status            string  `json:"status"`
	Receiver          string  `json:"receiver"`
	GroupLabels       Labels  `json:"groupLabels"`
	CommonLabels      Labels  `json:"commonLabels"`
	CommonAnnotations Labels  `json:"commonAnnotations"`
	ExternalURL       string  `json:"externalURL"`
	Alerts            []Alert `json:"alerts"`
}

type AlertStatus string

const (
	AlertStatusResolved AlertStatus = "resolved"
	AlertStatusFiring   AlertStatus = "firing"
)

type Alert struct {
	Status       AlertStatus `json:"status"`
	Labels       Labels      `json:"labels"`
	Annotations  Labels      `json:"annotations"`
	StartsAt     time.Time   `json:"startsAt"`
	EndsAt       time.Time   `json:"endsAt"`
	GeneratorURL string      `json:"generatorURL"`
	Fingerprint  string      `json:"fingerprint"`
}
