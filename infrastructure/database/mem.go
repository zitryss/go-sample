// Do not use this package in production! It is not optimized and was
// only developed to make testing without actual database possible.

package database

import (
	"bytes"
	"context"
	"encoding/gob"
	"sort"
	"time"

	. "github.com/zitryss/perfmon/domain"
	. "github.com/zitryss/perfmon/internal/context"
)

type mem struct {
	db []*record
}

type record struct {
	Job
	createdOn time.Time
	createdBy string
}

func NewMem() *mem {
	m := &mem{}
	m.populate()
	return m
}

func (m *mem) populate() {
	timestamps := []time.Time{
		time.Date(2017, 7, 1, 6, 0, 4, 0, time.UTC),
		time.Date(2017, 8, 19, 8, 40, 23, 0, time.UTC),
		time.Date(2017, 10, 5, 12, 57, 44, 0, time.UTC),
		time.Date(2017, 10, 29, 16, 18, 7, 0, time.UTC),
		time.Date(2018, 2, 8, 0, 51, 55, 0, time.UTC),
		time.Date(2018, 2, 16, 2, 16, 51, 0, time.UTC),
		time.Date(2018, 3, 24, 4, 25, 43, 0, time.UTC),
		time.Date(2018, 5, 25, 4, 52, 34, 0, time.UTC),
		time.Date(2018, 7, 2, 6, 34, 54, 0, time.UTC),
		time.Date(2018, 7, 11, 7, 26, 41, 0, time.UTC),
	}
	records := []*record{
		{Job{"p2", "v21", nil, "", "", time.Time{}, 0}, time.Now(), ""},
		{Job{"p2", "v22", nil, "", "", time.Time{}, 0}, time.Now(), ""},
		{Job{"p2", "v22", nil, "", "", time.Time{}, 0}, time.Now(), ""},
		{Job{"p3", "v3", []string{"a31", "a32"}, "", "", time.Time{}, 0}, time.Now(), ""},
		{Job{"p3", "v3", []string{"a32", "a33"}, "", "", time.Time{}, 0}, time.Now(), ""},
		{Job{"pY", "v3", []string{"a33", "a34"}, "", "", time.Time{}, 0}, time.Now(), ""},
		{Job{"p3", "vY", []string{"a34", "a35"}, "", "", time.Time{}, 0}, time.Now(), ""},
		{Job{"p4", "v4", []string{"a41", "a42", "a43"}, "n41", "", time.Time{}, 0}, time.Now(), ""},
		{Job{"p4", "v4", []string{"a41", "a42", "a43"}, "n42", "", time.Time{}, 0}, time.Now(), ""},
		{Job{"p4", "v4", []string{"a41", "a42", "a43"}, "n42", "", time.Time{}, 0}, time.Now(), ""},
		{Job{"pY", "v4", []string{"a41", "a42", "a43"}, "n43", "", time.Time{}, 0}, time.Now(), ""},
		{Job{"p4", "vY", []string{"a41", "a42", "a43"}, "n44", "", time.Time{}, 0}, time.Now(), ""},
		{Job{"p4", "v4", []string{"a41", "a4Y", "a43"}, "n45", "", time.Time{}, 0}, time.Now(), ""},
		{Job{"p5", "v5", []string{"a51", "a52", "a53"}, "n5", "m51", time.Time{}, 0}, time.Now(), ""},
		{Job{"p5", "v5", []string{"a51", "a52", "a53"}, "n5", "m52", time.Time{}, 0}, time.Now(), ""},
		{Job{"p5", "v5", []string{"a51", "a52", "a53"}, "n5", "m52", time.Time{}, 0}, time.Now(), ""},
		{Job{"pY", "v5", []string{"a51", "a52", "a53"}, "n5", "m53", time.Time{}, 0}, time.Now(), ""},
		{Job{"p5", "vY", []string{"a51", "a52", "a53"}, "n5", "m54", time.Time{}, 0}, time.Now(), ""},
		{Job{"p5", "v5", []string{"a51", "a5Y", "a53"}, "n5", "m55", time.Time{}, 0}, time.Now(), ""},
		{Job{"p5", "v5", []string{"a51", "a52", "a53"}, "nY", "m56", time.Time{}, 0}, time.Now(), ""},
		{Job{"p6", "v6", []string{"a61", "a62", "a63"}, "n6", "m6", timestamps[0], 27}, time.Now(), ""},
		{Job{"p6", "v6", []string{"a61", "a62", "a63"}, "n6", "m6", timestamps[1], 10}, time.Now(), ""},
		{Job{"p6", "v6", []string{"a61", "a62", "a63"}, "n6", "m6", timestamps[2], 33}, time.Now(), ""},
		{Job{"p6", "v6", []string{"a61", "a62", "a63"}, "n6", "m6", timestamps[3], 48}, time.Now(), ""},
		{Job{"pY", "v6", []string{"a61", "a62", "a63"}, "n6", "m6", timestamps[4], 74}, time.Now(), ""},
		{Job{"p6", "vY", []string{"a61", "a62", "a63"}, "n6", "m6", timestamps[5], 37}, time.Now(), ""},
		{Job{"p6", "v6", []string{"a61", "a6Y", "a63"}, "n6", "m6", timestamps[6], 92}, time.Now(), ""},
		{Job{"p6", "v6", []string{"a61", "a62", "a63"}, "nY", "m6", timestamps[7], 11}, time.Now(), ""},
		{Job{"p6", "v6", []string{"a61", "a62", "a63"}, "n6", "mY", timestamps[8], 42}, time.Now(), ""},
		{Job{"p6", "v6", []string{"a61", "a62", "a63"}, "n6", "m6", timestamps[9], 12}, time.Now(), ""},
		{Job{"p7", "v7", []string{"a71", "a72", "a73"}, "n7", "m7", timestamps[4], 0}, time.Now(), ""},
		{Job{"p7", "v7", []string{"a71", "a72", "a74"}, "n7", "m7", timestamps[5], 0}, time.Now(), ""},
	}
	m.db = records
}

func (m *mem) Create(ctx context.Context, job *Job) error {
	for _, r := range m.db {
		if job.Product != r.Product {
			continue
		}
		if job.Version != r.Version {
			continue
		}
		set1 := newSet(job.Attributes...)
		set2 := newSet(r.Attributes...)
		if !set1.isEqual(set2) {
			continue
		}
		if job.Name != r.Name {
			continue
		}
		if job.Measurement != r.Measurement {
			continue
		}
		if !job.Timestamp.Equal(r.Timestamp) {
			continue
		}
		if job.Value != r.Value {
			continue
		}
		return ErrJobAlreadyExists
	}
	r := &record{*job, CreatedOn(ctx), CreatedBy(ctx)}
	m.db = append(m.db, r)
	return nil
}

// Delete is not the part of the Databaser interface and only exists to
// cancel out the effect of the method Create in the database.
func (m *mem) Delete(job *Job) error {
	for i, r := range m.db {
		if job.Product != r.Product {
			continue
		}
		if job.Version != r.Version {
			continue
		}
		set1 := newSet(job.Attributes...)
		set2 := newSet(r.Attributes...)
		if !set1.isEqual(set2) {
			continue
		}
		if job.Name != r.Name {
			continue
		}
		if job.Measurement != r.Measurement {
			continue
		}
		if !job.Timestamp.Equal(r.Timestamp) {
			continue
		}
		if job.Value != r.Value {
			continue
		}
		m.db[i] = m.db[len(m.db)-1]
		m.db = m.db[:len(m.db)-1]
	}
	return nil
}

func (m *mem) Products() ([]string, error) {
	s := newSet()
	for _, r := range m.db {
		s.add(r.Product)
	}
	products := s.elements()
	return products, nil
}

func (m *mem) Versions(product string) ([]string, error) {
	s := newSet()
	for _, r := range m.db {
		if product != r.Product {
			continue
		}
		s.add(r.Version)
	}
	versions := s.elements()
	return versions, nil
}

func (m *mem) Attributes(product string, version string) ([]string, error) {
	s := newSet()
	for _, r := range m.db {
		if product != r.Product {
			continue
		}
		if version != r.Version {
			continue
		}
		s.add(r.Attributes...)
	}
	attributes := s.elements()
	return attributes, nil
}

func (m *mem) Names(product string, version string, attributes []string) ([]string, error) {
	s := newSet()
	for _, r := range m.db {
		if product != r.Product {
			continue
		}
		if version != r.Version {
			continue
		}
		set1 := newSet(attributes...)
		set2 := newSet(r.Attributes...)
		if !set1.isSubset(set2) {
			continue
		}
		s.add(r.Name)
	}
	names := s.elements()
	return names, nil
}

func (m *mem) Measurements(product string, version string, attributes []string, name string) ([]string, error) {
	s := newSet()
	for _, r := range m.db {
		if product != r.Product {
			continue
		}
		if version != r.Version {
			continue
		}
		set1 := newSet(attributes...)
		set2 := newSet(r.Attributes...)
		if !set1.isSubset(set2) {
			continue
		}
		if name != r.Name {
			continue
		}
		s.add(r.Measurement)
	}
	measurements := s.elements()
	return measurements, nil
}

func (m *mem) Jobs(product string, version string, attributes []string, name string, measurement string, lbound time.Time, rbound time.Time) ([]*Job, error) {
	s := newSet()
	for _, r := range m.db {
		if product != r.Product {
			continue
		}
		if version != r.Version {
			continue
		}
		set1 := newSet(attributes...)
		set2 := newSet(r.Attributes...)
		if !set1.isSubset(set2) {
			continue
		}
		if name != r.Name {
			continue
		}
		if measurement != r.Measurement {
			continue
		}
		if r.Timestamp.Before(lbound) || r.Timestamp.After(rbound) {
			continue
		}
		buf := &bytes.Buffer{}
		enc := gob.NewEncoder(buf)
		err := enc.Encode(r.Job)
		if err != nil {
			return nil, err
		}
		s.add(buf.String())
	}
	ee := s.elements()
	jobs := ([]*Job)(nil)
	for _, e := range ee {
		buf := &bytes.Buffer{}
		dec := gob.NewDecoder(buf)
		buf.WriteString(e)
		j := &Job{}
		err := dec.Decode(j)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

type set map[string]struct{}

func newSet(elems ...string) set {
	s := set(map[string]struct{}{})
	s.add(elems...)
	return s
}

func (s set) add(elems ...string) {
	for _, e := range elems {
		s[e] = struct{}{}
	}
}

func (s set) elements() []string {
	elems := ([]string)(nil)
	for k := range s {
		elems = append(elems, k)
	}
	sort.Strings(elems)
	return elems
}

func (s set) isEqual(other set) bool {
	if len(s) != len(other) {
		return false
	}
	if !s.isSubset(other) {
		return false
	}
	return true
}

func (s set) isSubset(other set) bool {
	for k := range s {
		_, ok := other[k]
		if !ok {
			return false
		}
	}
	return true
}
