package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotEnoughAttributes = errors.New("not enough attributes provided")
	ErrJobAlreadyExists    = errors.New("job already exists")
)

type Job struct {
	Product     string
	Version     string
	Attributes  []string
	Name        string
	Measurement string
	Timestamp   time.Time
	Value       int
}

type Usecaser interface {
	Add(ctx context.Context, job *Job) error
	Products() ([]string, error)
	Versions(product string) ([]string, error)
	Attributes(product string, version string) ([]string, error)
	Names(product string, version string, attributes []string) ([]string, error)
	Measurements(product string, version string, attributes []string, name string) ([]string, error)
	Jobs(product string, version string, attributes []string, name string, measurement string, lbound time.Time, rbound time.Time) ([]*Job, error)
}

type Databaser interface {
	Create(ctx context.Context, job *Job) error
	Products() ([]string, error)
	Versions(product string) ([]string, error)
	Attributes(product string, version string) ([]string, error)
	Names(product string, version string, attributes []string) ([]string, error)
	Measurements(product string, version string, attributes []string, name string) ([]string, error)
	Jobs(product string, version string, attributes []string, name string, measurement string, lbound time.Time, rbound time.Time) ([]*Job, error)
}
