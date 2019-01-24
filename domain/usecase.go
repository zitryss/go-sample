package domain

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type usecase struct {
	db Databaser
}

func NewUsecase(db Databaser) *usecase {
	return &usecase{db}
}

func (u *usecase) Add(ctx context.Context, job *Job) error {
	err := u.db.Create(ctx, job)
	if err != nil {
		return err
	}
	return nil
}

func (u *usecase) Products() ([]string, error) {
	products, err := u.db.Products()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (u *usecase) Versions(product string) ([]string, error) {
	versions, err := u.db.Versions(product)
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func (u *usecase) Attributes(product string, version string) ([]string, error) {
	attributes, err := u.db.Attributes(product, version)
	if err != nil {
		return nil, err
	}
	return attributes, nil
}

func (u *usecase) Names(product string, version string, attributes []string) ([]string, error) {
	names, err := u.db.Names(product, version, attributes)
	if err != nil {
		return nil, err
	}
	return names, nil
}

func (u *usecase) Measurements(product string, version string, attributes []string, name string) ([]string, error) {
	measurements, err := u.db.Measurements(product, version, attributes, name)
	if err != nil {
		return nil, err
	}
	return measurements, nil
}

func (u *usecase) Jobs(product string, version string, attributes []string, name string, measurement string, lbound time.Time, rbound time.Time) ([]*Job, error) {
	jobs, err := u.db.Jobs(product, version, attributes, name, measurement, lbound, rbound)
	if err != nil {
		return nil, err
	}
	err = checkAttributes(jobs)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// checkAttributes makes sure that all the jobs have the same set of
// attributes.
func checkAttributes(jobs []*Job) error {
	s := map[string]struct{}{}
	for _, j := range jobs {
		sort.Strings(j.Attributes)
		attributes := strings.Join(j.Attributes, "")
		s[attributes] = struct{}{}
	}
	if len(s) > 1 {
		err := ErrNotEnoughAttributes
		err = errors.WithStack(err)
		return err
	}
	return nil
}
