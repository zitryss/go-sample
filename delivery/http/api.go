package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	. "github.com/zitryss/perfmon/domain"
	. "github.com/zitryss/perfmon/internal/context"
	"github.com/zitryss/perfmon/internal/log"
)

type api struct {
	use Usecaser
}

func newAPI(use Usecaser) *api {
	return &api{use}
}

func (a *api) upload() httprouter.Handle {
	type request struct {
		Product     string
		Version     string
		Attributes  []string
		Name        string
		Measurement string
		Timestamp   time.Time
		Value       int
	}
	input := func(r *http.Request) (context.Context, *request, error) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			err = errors.WithStack(err)
			return nil, nil, err
		}
		if req.Product == "" {
			return nil, nil, errors.New("product is missing")
		}
		if req.Version == "" {
			return nil, nil, errors.New("version is missing")
		}
		if req.Name == "" {
			return nil, nil, errors.New("name is missing")
		}
		if req.Measurement == "" {
			return nil, nil, errors.New("measurement is missing")
		}
		if req.Timestamp.IsZero() {
			return nil, nil, errors.New("timestamp is missing")
		}
		if req.Value == 0 {
			return nil, nil, errors.New("value is missing")
		}
		ctx := r.Context()
		ctx = WithCreatedOn(ctx, time.Now())
		ctx = WithCreatedBy(ctx, r.RemoteAddr)
		return ctx, req, nil
	}
	process := func(ctx context.Context, req *request) error {
		job := &Job{}
		job.Product = req.Product
		job.Version = req.Version
		job.Attributes = req.Attributes
		job.Name = req.Name
		job.Measurement = req.Measurement
		job.Timestamp = req.Timestamp
		job.Value = req.Value
		err := a.use.Add(ctx, job)
		if err != nil {
			return err
		}
		return nil
	}
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx, req, err := input(r)
		if err != nil {
			err = errors.WithMessage(err, "bad request")
			log.Debug(err)
			http.Error(w, "Bad Request", 400)
			return
		}
		err = process(ctx, req)
		if errors.Cause(err) == ErrJobAlreadyExists {
			err = errors.WithMessage(err, "bad request")
			log.Debug(err)
			http.Error(w, "Resource Already Exists", 409)
			return
		}
		if err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		w.WriteHeader(201)
	}
}

func (a *api) products() httprouter.Handle {
	type response struct {
		Products []string `json:"products"`
	}
	process := func() ([]string, error) {
		data, err := a.use.Products()
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	output := func(w http.ResponseWriter, data []string) error {
		resp := &response{}
		resp.Products = data
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			err = errors.WithStack(err)
			return err
		}
		return nil
	}
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		data, err := process()
		if err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		err = output(w, data)
		if err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}
}

func (a *api) versions() httprouter.Handle {
	type request struct {
		Product string
	}
	type response struct {
		Versions []string `json:"versions"`
	}
	input := func(r *http.Request) (*request, error) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		if req.Product == "" {
			return nil, errors.New("product is missing")
		}
		return req, nil
	}
	process := func(req *request) ([]string, error) {
		data, err := a.use.Versions(req.Product)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	output := func(w http.ResponseWriter, data []string) error {
		resp := &response{}
		resp.Versions = data
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			err = errors.WithStack(err)
			return err
		}
		return nil
	}
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		req, err := input(r)
		if err != nil {
			err = errors.WithMessage(err, "bad request")
			log.Debug(err)
			http.Error(w, "Bad Request", 400)
			return
		}
		data, err := process(req)
		if err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		err = output(w, data)
		if err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}
}

func (a *api) attributes() httprouter.Handle {
	type request struct {
		Product string
		Version string
	}
	type response struct {
		Attributes []string `json:"attributes"`
	}
	input := func(r *http.Request) (*request, error) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		if req.Product == "" {
			return nil, errors.New("product is missing")
		}
		if req.Version == "" {
			return nil, errors.New("version is missing")
		}
		return req, nil
	}
	process := func(req *request) ([]string, error) {
		data, err := a.use.Attributes(req.Product, req.Version)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	output := func(w http.ResponseWriter, data []string) error {
		resp := &response{}
		resp.Attributes = data
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			err = errors.WithStack(err)
			return err
		}
		return nil
	}
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		req, err := input(r)
		if err != nil {
			err = errors.WithMessage(err, "bad request")
			log.Debug(err)
			http.Error(w, "Bad Request", 400)
			return
		}
		data, err := process(req)
		if err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		err = output(w, data)
		if err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}
}

func (a *api) names() httprouter.Handle {
	type request struct {
		Product    string
		Version    string
		Attributes []string
	}
	type response struct {
		Names []string `json:"names"`
	}
	input := func(r *http.Request) (*request, error) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		if req.Product == "" {
			return nil, errors.New("product is missing")
		}
		if req.Version == "" {
			return nil, errors.New("version is missing")
		}
		return req, nil
	}
	process := func(req *request) ([]string, error) {
		data, err := a.use.Names(req.Product, req.Version, req.Attributes)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	output := func(w http.ResponseWriter, data []string) error {
		resp := &response{}
		resp.Names = data
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			err = errors.WithStack(err)
			return err
		}
		return nil
	}
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		req, err := input(r)
		if err != nil {
			err = errors.WithMessage(err, "bad request")
			log.Debug(err)
			http.Error(w, "Bad Request", 400)
			return
		}
		data, err := process(req)
		if err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		err = output(w, data)
		if err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}
}

func (a *api) measurements() httprouter.Handle {
	type request struct {
		Product    string
		Version    string
		Attributes []string
		Name       string
	}
	type response struct {
		Measurements []string `json:"measurements"`
	}
	input := func(r *http.Request) (*request, error) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		if req.Product == "" {
			return nil, errors.New("product is missing")
		}
		if req.Version == "" {
			return nil, errors.New("version is missing")
		}
		if req.Name == "" {
			return nil, errors.New("name is missing")
		}
		return req, nil
	}
	process := func(req *request) ([]string, error) {
		data, err := a.use.Measurements(req.Product, req.Version, req.Attributes, req.Name)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	output := func(w http.ResponseWriter, data []string) error {
		resp := &response{}
		resp.Measurements = data
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			err = errors.WithStack(err)
			return err
		}
		return nil
	}
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		req, err := input(r)
		if err != nil {
			err = errors.WithMessage(err, "bad request")
			log.Debug(err)
			http.Error(w, "Bad Request", 400)
			return
		}
		data, err := process(req)
		if err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		err = output(w, data)
		if err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}
}

func (a *api) chart() httprouter.Handle {
	type request struct {
		Product     string
		Version     string
		Attributes  []string
		Name        string
		Measurement string
		LBound      time.Time
		RBound      time.Time
	}
	type Point struct {
		X time.Time `json:"x"`
		Y int       `json:"y"`
	}
	type response struct {
		Data []*Point `json:"data"`
	}
	input := func(r *http.Request) (*request, error) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		if req.Product == "" {
			return nil, errors.New("product is missing")
		}
		if req.Version == "" {
			return nil, errors.New("version is missing")
		}
		if req.Name == "" {
			return nil, errors.New("name is missing")
		}
		if req.Measurement == "" {
			return nil, errors.New("measurement is missing")
		}
		if req.LBound.IsZero() {
			return nil, errors.New("left bound is missing")
		}
		if req.RBound.IsZero() {
			return nil, errors.New("right bound is missing")
		}
		if req.LBound.After(req.RBound) {
			return nil, errors.New("left bound should be smaller than the right bound")
		}
		return req, nil
	}
	process := func(req *request) ([]*Job, error) {
		data, err := a.use.Jobs(req.Product, req.Version, req.Attributes, req.Name, req.Measurement, req.LBound, req.RBound)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	output := func(w http.ResponseWriter, data []*Job) error {
		resp := &response{}
		for _, d := range data {
			p := &Point{d.Timestamp, d.Value}
			resp.Data = append(resp.Data, p)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			err = errors.WithStack(err)
			return err
		}
		return nil
	}
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		req, err := input(r)
		if err != nil {
			err = errors.WithMessage(err, "bad request")
			log.Debug(err)
			http.Error(w, "Bad Request", 400)
			return
		}
		data, err := process(req)
		if errors.Cause(err) == ErrNotEnoughAttributes {
			err = errors.WithMessage(err, "bad request")
			log.Debug(err)
			http.Error(w, "Not enough attributes provided!", 400)
			return
		}
		if err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		err = output(w, data)
		if err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}
}
