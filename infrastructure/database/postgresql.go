package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	. "github.com/zitryss/perfmon/domain"
	. "github.com/zitryss/perfmon/internal/context"
	"github.com/zitryss/perfmon/internal/log"
)

type psql struct {
	db *sql.DB
}

func NewPsql(schema string) *psql {
	connStr := fmt.Sprintf("postgresql://postgres:mysecretpassword@db:5432/postgres?search_path=%s&sslmode=disable", schema)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Critical("database: postgresql: open")
		log.Critical(err)
		os.Exit(1)
	}
	err = db.Ping()
	if err != nil {
		log.Critical("database: postgresql: ping")
		log.Critical(err)
		os.Exit(1)
	}
	return &psql{db}
}

func (p *psql) Create(ctx context.Context, job *Job) error {
	jIDsFM, err := p.findJobsFullMatch(job)
	if err != nil {
		return err
	}
	if len(jIDsFM) > 0 {
		err := ErrJobAlreadyExists
		err = errors.WithStack(err)
		return err
	}
	jID, err := p.createJob(ctx, job.Product, job.Version, job.Name, job.Measurement, job.Timestamp, job.Value)
	if err != nil {
		return err
	}
	aIDs, err := p.createAttributes(job.Attributes)
	if err != nil {
		return err
	}
	err = p.createJobAttributes(jID, aIDs)
	if err != nil {
		return err
	}
	return nil
}

func (p *psql) createJob(ctx context.Context, product string, version string, name string, measurement string, timestamp time.Time, value int) (string, error) {
	query := `
		INSERT INTO job (product, version, name, measurement, timestamp, value, created_on, created_by)
		     VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		  RETURNING id
	`
	jID := ""
	err := p.db.QueryRow(query, product, version, name, measurement, timestamp, value, CreatedOn(ctx), CreatedBy(ctx)).Scan(&jID)
	if err != nil {
		err = errors.WithStack(err)
		return "", err
	}
	return jID, nil
}

func (p *psql) createAttributes(attributes []string) ([]string, error) {
	aIDs := ([]string)(nil)
	for _, a := range attributes {
		aID, found, err := p.findAttribute(a)
		if err != nil {
			return nil, err
		}
		if found {
			aIDs = append(aIDs, aID)
			continue
		}
		aID, err = p.createAttribute(a)
		if err != nil {
			return nil, err
		}
		aIDs = append(aIDs, aID)
	}
	return aIDs, nil
}

func (p *psql) findAttribute(name string) (string, bool, error) {
	query := `
		SELECT a.id
		  FROM attribute AS a
		 WHERE a.name = $1
	`
	aID := ""
	err := p.db.QueryRow(query, name).Scan(&aID)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		err = errors.WithStack(err)
		return "", false, err
	}
	return aID, true, nil
}

func (p *psql) createAttribute(name string) (string, error) {
	query := `
		INSERT INTO attribute (name)
		     VALUES ($1)
		  RETURNING id
	`
	aID := ""
	err := p.db.QueryRow(query, name).Scan(&aID)
	if err != nil {
		err = errors.WithStack(err)
		return "", err
	}
	return aID, nil
}

func (p *psql) createJobAttributes(jID string, aIDs []string) error {
	for _, aID := range aIDs {
		_, err := p.createJobAttribute(jID, aID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *psql) createJobAttribute(jID string, aID string) (string, error) {
	query := `
		INSERT INTO job_attribute (job_id, attribute_id)
		     VALUES ($1, $2)
		  RETURNING id
	`
	jaID := ""
	err := p.db.QueryRow(query, jID, aID).Scan(&jaID)
	if err != nil {
		err = errors.WithStack(err)
		return "", err
	}
	return jaID, nil
}

// Delete is not the part of the Databaser interface and only exists to
// cancel out the effect of the method Create in the database.
func (p *psql) Delete(job *Job) error {
	jIDsFM, err := p.findJobsFullMatch(job)
	if err != nil {
		return err
	}
	if len(jIDsFM) == 0 {
		return nil
	}
	err = p.deleteJobAttributes(jIDsFM)
	if err != nil {
		return err
	}
	err = p.deleteAttributes()
	if err != nil {
		return err
	}
	err = p.deleteJob(jIDsFM)
	if err != nil {
		return err
	}
	return nil
}

func (p *psql) findJobsFullMatch(job *Job) ([]string, error) {
	jIDs, err := p.findJobs(job.Product, job.Version, job.Name, job.Measurement, job.Timestamp, job.Value)
	if err != nil {
		return nil, err
	}
	jIDsFM := ([]string)(nil)
	for _, jID := range jIDs {
		attributes, err := p.attributes(jID)
		if err != nil {
			return nil, err
		}
		if equalSlice(job.Attributes, attributes) {
			jIDsFM = append(jIDsFM, jID)
		}
	}
	return jIDsFM, nil
}

func (p *psql) findJobs(product string, version string, name string, measurement string, timestamp time.Time, value int) ([]string, error) {
	query := `
		SELECT j.id
		  FROM job AS j
		 WHERE j.product = $1
		       AND j.version = $2
		       AND j.name = $3
		       AND j.measurement = $4
		       AND j.timestamp = $5
		       AND j.value = $6
	`
	rows, err := p.db.Query(query, product, version, name, measurement, timestamp, value)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	defer rows.Close()
	jIDs := ([]string)(nil)
	jID := ""
	for rows.Next() {
		err = rows.Scan(&jID)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		jIDs = append(jIDs, jID)
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return jIDs, nil
}

func equalSlice(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func (p *psql) deleteJobAttributes(jIDs []string) error {
	query := `
		DELETE FROM job_attribute AS ja
		      WHERE ja.job_id = ANY($1)
	`
	_, err := p.db.Exec(query, pq.StringArray(jIDs))
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	return nil
}

func (p *psql) deleteAttributes() error {
	query := `
		DELETE FROM attribute AS a
		      WHERE a.id NOT IN (
		            	SELECT a.id
		            	  FROM job_attribute AS ja
		            	  JOIN attribute AS a
		            	       ON ja.attribute_id = a.id
		            )
	`
	_, err := p.db.Exec(query)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	return nil
}

func (p *psql) deleteJob(jIDs []string) error {
	query := `
		DELETE FROM job AS j
		      WHERE j.id = ANY($1)
	`
	_, err := p.db.Exec(query, pq.StringArray(jIDs))
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	return nil
}

func (p *psql) Products() ([]string, error) {
	query := `
		  SELECT DISTINCT j.product
		    FROM job AS j
		ORDER BY j.product
	`
	rows, err := p.db.Query(query)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	defer rows.Close()
	products := ([]string)(nil)
	prod := ""
	for rows.Next() {
		err = rows.Scan(&prod)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		products = append(products, prod)
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return products, nil
}

func (p *psql) Versions(product string) ([]string, error) {
	query := `
		  SELECT DISTINCT j.version
		    FROM job AS j
		   WHERE j.product = $1
		ORDER BY j.version
	`
	rows, err := p.db.Query(query, product)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	defer rows.Close()
	versions := ([]string)(nil)
	v := ""
	for rows.Next() {
		err = rows.Scan(&v)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		versions = append(versions, v)
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return versions, nil
}

func (p *psql) Attributes(product string, version string) ([]string, error) {
	query := `
		  SELECT DISTINCT a.name
		    FROM job AS j
		    JOIN job_attribute AS ja
		         ON j.id = ja.job_id
		    JOIN attribute AS a
		         ON ja.attribute_id = a.id
		   WHERE j.product = $1
		         AND j.version = $2
		ORDER BY a.name
	`
	rows, err := p.db.Query(query, product, version)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	defer rows.Close()
	attributes := ([]string)(nil)
	a := ""
	for rows.Next() {
		err = rows.Scan(&a)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		attributes = append(attributes, a)
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return attributes, nil
}

func (p *psql) Names(product string, version string, attributes []string) ([]string, error) {
	query := `
		  SELECT DISTINCT j.name
		    FROM job AS j
		    JOIN job_attribute AS ja
		         ON j.id = ja.job_id
		    JOIN attribute AS a
		         ON ja.attribute_id = a.id
		   WHERE j.product = $1
		         AND j.version = $2
		         AND a.name = ANY($3)
		GROUP BY j.id
		  HAVING COUNT(j.id) = ARRAY_LENGTH($3, 1)
		ORDER BY j.name
	`
	rows, err := p.db.Query(query, product, version, pq.StringArray(attributes))
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	defer rows.Close()
	names := ([]string)(nil)
	n := ""
	for rows.Next() {
		err = rows.Scan(&n)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		names = append(names, n)
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return names, nil
}

func (p *psql) Measurements(product string, version string, attributes []string, name string) ([]string, error) {
	query := `
		  SELECT DISTINCT j.measurement
		    FROM job AS j
		    JOIN job_attribute AS ja
		         ON j.id = ja.job_id
		    JOIN attribute AS a
		         ON ja.attribute_id = a.id
		   WHERE j.product = $1
		         AND j.version = $2
		         AND a.name = ANY($3)
		         AND j.name = $4
		GROUP BY j.id
		  HAVING COUNT(j.id) = ARRAY_LENGTH($3, 1)
		ORDER BY j.measurement
	`
	rows, err := p.db.Query(query, product, version, pq.StringArray(attributes), name)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	defer rows.Close()
	measurements := ([]string)(nil)
	m := ""
	for rows.Next() {
		err = rows.Scan(&m)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		measurements = append(measurements, m)
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return measurements, nil
}

func (p *psql) Jobs(product string, version string, attributes []string, name string, measurement string, lbound time.Time, rbound time.Time) ([]*Job, error) {
	query := `
		  SELECT j.id, j.product, j.version, j.name, j.measurement, j.timestamp, j.value
		    FROM job AS j
		    JOIN job_attribute AS ja
		         ON j.id = ja.job_id
		    JOIN attribute AS a
		         ON ja.attribute_id = a.id
		   WHERE j.product = $1
		         AND j.version = $2
		         AND a.name = ANY($3)
		         AND j.name = $4
		         AND j.measurement = $5
		         AND j.timestamp BETWEEN $6 AND $7
		GROUP BY j.id
		  HAVING COUNT(j.id) = ARRAY_LENGTH($3, 1)
		ORDER BY j.timestamp
	`
	rows, err := p.db.Query(query, product, version, pq.StringArray(attributes), name, measurement, lbound, rbound)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	defer rows.Close()
	jobs := ([]*Job)(nil)
	jID := ""
	for rows.Next() {
		j := &Job{}
		err = rows.Scan(&jID, &j.Product, &j.Version, &j.Name, &j.Measurement, &j.Timestamp, &j.Value)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		attributes, err = p.attributes(jID)
		if err != nil {
			return nil, err
		}
		j.Attributes = attributes
		jobs = append(jobs, j)
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return jobs, nil
}

func (p *psql) attributes(jID string) ([]string, error) {
	query := `
		SELECT a.name
		  FROM job_attribute AS ja
		  JOIN attribute AS a
		       ON ja.attribute_id = a.id
		 WHERE ja.job_id = $1
	`
	rows, err := p.db.Query(query, jID)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	defer rows.Close()
	attributes := ([]string)(nil)
	a := ""
	for rows.Next() {
		err = rows.Scan(&a)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		attributes = append(attributes, a)
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return attributes, nil
}
