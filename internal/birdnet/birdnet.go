package birdnet

import (
	"fmt"
	"time"
)

// BirdImage represents the bird image metadata
type BirdImage struct {
	URL            string    `json:"URL"`
	ScientificName string    `json:"ScientificName"`
	LicenseName    string    `json:"LicenseName"`
	LicenseURL     string    `json:"LicenseURL"`
	AuthorName     string    `json:"AuthorName"`
	AuthorURL      string    `json:"AuthorURL"`
	CachedAt       time.Time `json:"CachedAt"`
	SourceProvider string    `json:"SourceProvider"`
}

// BirdDetection represents the main bird detection data
type BirdDetection struct {
	ID             int        `json:"ID"`
	SourceNode     string     `json:"SourceNode"`
	Date           string     `json:"Date"`
	Time           string     `json:"Time"`
	Source         string     `json:"Source"`
	BeginTime      time.Time  `json:"BeginTime"`
	EndTime        time.Time  `json:"EndTime"`
	SpeciesCode    string     `json:"SpeciesCode"`
	ScientificName string     `json:"ScientificName"`
	CommonName     string     `json:"CommonName"`
	Confidence     float64    `json:"Confidence"`
	Latitude       float64    `json:"Latitude"`
	Longitude      float64    `json:"Longitude"`
	Threshold      float64    `json:"Threshold"`
	Sensitivity    float64    `json:"Sensitivity"`
	ClipName       string     `json:"ClipName"`
	ProcessingTime int64      `json:"ProcessingTime"`
	Results        *string    `json:"Results"`
	Review         *string    `json:"Review"`
	Comments       *string    `json:"Comments"`
	Lock           *string    `json:"Lock"`
	Verified       string     `json:"Verified"`
	Locked         bool       `json:"Locked"`
	BirdImage      *BirdImage `json:"BirdImage"`
}

type BirdDetectionEvent struct {
	Time           time.Time `json:"Time"`
	SourceNode     string    `json:"SourceNode"`
	Source         string    `json:"Source"`
	BeginTime      time.Time `json:"BeginTime"`
	EndTime        time.Time `json:"EndTime"`
	SpeciesCode    string    `json:"SpeciesCode"`
	ScientificName string    `json:"ScientificName"`
	CommonName     string    `json:"CommonName"`
	Confidence     float64   `json:"Confidence"`
	Latitude       float64   `json:"Latitude"`
	Longitude      float64   `json:"Longitude"`
	Threshold      float64   `json:"Threshold"`
	Sensitivity    float64   `json:"Sensitivity"`
}

func (b *BirdDetection) ToBirdDetectionEvent(timezone string) (*BirdDetectionEvent, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("failed to load timezone %s: %w", timezone, err)
	}
	
	timedate := fmt.Sprintf("%s %s", b.Date, b.Time)
	eventTime, err := time.ParseInLocation("2006-01-02 15:04:05", timedate, loc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %w", err)
	}

	return &BirdDetectionEvent{
		SourceNode:     b.SourceNode,
		Time:           eventTime,
		Source:         b.Source,
		BeginTime:      b.BeginTime,
		EndTime:        b.EndTime,
		SpeciesCode:    b.SpeciesCode,
		ScientificName: b.ScientificName,
		CommonName:     b.CommonName,
		Confidence:     b.Confidence,
		Latitude:       b.Latitude,
		Longitude:      b.Longitude,
		Threshold:      b.Threshold,
		Sensitivity:    b.Sensitivity,
	}, err
}
