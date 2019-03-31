// Copyright 2019 Adverax. All Rights Reserved.
// This file is part of project
//
//      http://github.com/adverax/echo
//
// Licensed under the MIT (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://github.com/adverax/echo/blob/master/LICENSE
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package widget

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/adverax/echo"
	"github.com/adverax/echo/generic"
	"github.com/adverax/echo/security"
)

type MultiStepState struct {
	Stage      string         `json:"stage"`      // Stage name
	History    []string       `json:"history"`    // List of stage names
	Params     generic.Params `json:"values"`     // Map of accepted parameters
	Expiration int64          `json:"expiration"` // Expiration (UNIX timestamp)
}

// Become to the new state
// newState is variant:
//  * string - name of new state
//  * numeric - relative offset for reduce
func (s *MultiStepState) Become(
	newState interface{},
) error {
	switch stg := newState.(type) {
	case string:
		s.History = append(s.History, s.Stage)
		s.Stage = stg
		return nil
	case int:
		return s.reduce(stg)
	case int8:
		return s.reduce(int(stg))
	case int16:
		return s.reduce(int(stg))
	case int32:
		return s.reduce(int(stg))
	case int64:
		return s.reduce(int(stg))
	case uint:
		return s.reduce(int(stg))
	case uint8:
		return s.reduce(int(stg))
	case uint16:
		return s.reduce(int(stg))
	case uint32:
		return s.reduce(int(stg))
	case uint64:
		return s.reduce(int(stg))
	default:
		return errors.New("unknown stage type")
	}
}

// Reduce state by history
func (s *MultiStepState) reduce(
	offset int,
) error {
	if offset == 0 {
		return nil
	}

	index := len(s.History) + offset
	if index < 0 {
		return errors.New("range check error")
	}

	s.Stage = s.History[index]
	s.History = s.History[:index]
	return nil
}

type MultiStepStage interface {
	// Get model of stage
	Model(ctx echo.Context, state *MultiStepState) (echo.Model, error)
	// Publish model
	Publish(ctx echo.Context, state *MultiStepState, model echo.Model) error
	// Consume stage data. Return new state reference or nil (see method MultiStepState.Become)
	Consume(ctx echo.Context, state *MultiStepState, data interface{}) (reply interface{}, err error)
}

var MultiStepPrevBtn = &Action{
	Label: MessageMultistepFormPrev,
}

var MultiStepNextBtn = &Action{
	Label: MessageMultistepFormNext,
}

// Prototype for all MultiStepStages
type MultiStepBaseStage struct {
	PrevBtn *Action
	NextBtn *Action
}

func (stage *MultiStepBaseStage) Publish(
	ctx echo.Context,
	state *MultiStepState,
	model echo.Model,
) error {
	if len(state.History) != 0 && stage.PrevBtn != nil {
		model["Prev"] = &Action{
			Label: stage.PrevBtn.Label,
			Action: ctx.Echo().UrlLinker.Collapse(
				ctx,
				ctx.Request().URL.Path+"/undo",
			),
		}
	}

	if stage.NextBtn != nil {
		submit := &FormSubmit{
			Name:  "action",
			Label: stage.NextBtn.Label,
			Items: redoDataSet,
		}
		err := submit.SetValue(ctx, []string{"redo"})
		if err != nil {
			return err
		}
		model["Next"] = submit
	}

	return nil
}

func (stage *MultiStepBaseStage) Consume(
	ctx echo.Context,
	state *MultiStepState,
	data interface{},
) (reply interface{}, err error) {
	return nil, nil
}

type MultiStepStageResource struct {
	MultiStepBaseStage
	Name    string
	Content echo.Widget
}

func (stage *MultiStepStageResource) Model(
	ctx echo.Context,
	state *MultiStepState,
) (echo.Model, error) {
	accepted := &FormHidden{
		Name:     stage.Name,
		Required: true,
	}

	return echo.Model{
		"Accepted": accepted,
	}, nil
}

func (stage *MultiStepStageResource) Consume(
	ctx echo.Context,
	state *MultiStepState,
	data interface{},
) (reply interface{}, err error) {
	return nil, nil
}

func (stage *MultiStepStageResource) Publish(
	ctx echo.Context,
	state *MultiStepState,
	model echo.Model,
) error {
	err := stage.MultiStepBaseStage.Publish(ctx, state, model)
	if err != nil {
		return err
	}

	content, err := echo.RenderWidget(ctx, stage.Content)
	if err != nil {
		return err
	}
	model["Content"] = content

	return nil
}

type MultiStepStrategy interface {
	// Started initialization. Executed for restart strategy only.
	Setup(ctx echo.Context, values generic.Params) (stage string, err error)
	// Create stage instance
	Stage(ctx echo.Context, state *MultiStepState) (MultiStepStage, error)
}

// Multi step form (with standard buttons)
// Example:
// func actionXXX(){
//   form := &widget.MultiStepForm{
//     Strategy: MyStrategy,
//     Security: security,
//     Timeout: 100,
//     ...
//   }
//   return form.Execute(env)
// }
//
// Define stages for the form (must implements MultiStepStage)
// type MyStage struct {
//    widgets.MultiStepBaseStage
// }
//
// Initialize router
// func Setup(){
//    InitMultiStepRouter(router, actionXXX)
// }
type MultiStepForm struct {
	Form
	Security GuidMaker    // Optional
	Storage  echo.Storage // Optional
	Strategy MultiStepStrategy
	Timeout  time.Duration
}

func (w *MultiStepForm) Execute(
	ctx echo.Context,
) (state *MultiStepState, err error) {
	state = new(MultiStepState)
	storage := w.getStorage(ctx)
	// Load form identifier
	request := ctx.Request()
	id := ctx.Param("id")
	if request.Method == http.MethodGet && id == "" {
		return nil, w.restart(ctx, storage)
	}

	key := makeMultiStepKey(id)

	// Restore current state
	err = storage.Get(key, &state)
	if err != nil {
		return nil, err
	}

	// Handle undo
	if ctx.Param("undo") == "undo" {
		err = state.reduce(-1)
		if err != nil {
			return nil, err
		}

		timeout := makeTimeout(state.Expiration)
		if timeout < 0 {
			err := storage.Delete(key)
			if err != nil {
				return nil, err
			}

			return state, w.redirect(ctx, id, 1)
		} else {
			err = storage.Set(key, state, timeout)
			if err != nil {
				return nil, err
			}

			return state, w.redirect(ctx, id, 1)
		}
	}

	// Create model
	stage, err := w.Strategy.Stage(ctx, state)
	if err != nil {
		_ = storage.Delete(key)
		return
	}

	model, err := stage.Model(ctx, state)
	if err != nil {
		_ = storage.Delete(key)
		return
	}
	if model == nil {
		return
	}

	// Assign model values
	for _, item := range model {
		if field, ok := item.(echo.ModelField); ok {
			name := field.GetName()
			if val, ok := state.Params[name]; ok {
				field.SetVal(ctx, val)
			}
		}
	}

	// Load and validate date
	if request.Method == http.MethodPost {
		err = model.Bind(ctx)
		if err != nil {
			return
		}
		if model.IsValid() {
			// Store state
			for _, item := range model {
				if field, ok := item.(echo.ModelField); ok {
					name := field.GetName()
					val := field.GetVal()
					if val != nil {
						state.Params[name] = val
					}
				}
			}

			reply, err := stage.Consume(ctx, state, model)
			if err != nil {
				return nil, err
			}
			if reply == nil {
				_ = storage.Delete(key)
				return nil, nil
			}

			if state.Expiration == 0 {
				state.Expiration = time.Now().Add(w.Timeout).Unix()
			}
			if state.Expiration == 0 {
				state.Expiration = TimeInfinity
			}
			timeout := makeTimeout(state.Expiration)
			if timeout < 0 {
				err := storage.Delete(key)
				if err != nil {
					return nil, err
				}

				return state, w.redirect(ctx, id, 0)
			}

			err = state.Become(reply)
			if err != nil {
				return nil, err
			}

			err = storage.Set(key, state, timeout)
			if err != nil {
				return nil, err
			}

			return state, w.redirect(ctx, id, 0)
		}
	}

	// Publish model
	w.Model = model
	return state, stage.Publish(ctx, state, model)
}

func (w *MultiStepForm) restart(
	ctx echo.Context,
	storage echo.Storage,
) error {
	id := strconv.FormatUint(w.getSecurity(ctx).CreateGuid(), 10)

	values := make(generic.Params, 16)
	stage, err := w.Strategy.Setup(ctx, values)
	if err != nil || stage == "" {
		return err
	}

	state := &MultiStepState{
		Stage:      stage,
		Params:     values,
		History:    make([]string, 0, 8),
		Expiration: time.Now().Add(w.Timeout).Unix(),
	}

	timeout := makeTimeout(state.Expiration)
	if timeout >= 0 {
		err = storage.Set(makeMultiStepKey(id), state, timeout)
		if err != nil {
			return err
		}
	}

	return w.redirect(ctx, id, 0)
}

func (w *MultiStepForm) redirect(
	ctx echo.Context,
	id string,
	skip int,
) error {
	req := ctx.Request()
	url := ctx.Echo().UrlLinker.Expand(ctx, req.URL.Path)
	slices := strings.Split(url, "/")
	if skip != 0 {
		slices = slices[:len(slices)-skip]
	}
	length := len(id)
	for i := len(slices) - 1; i >= 0; i-- {
		slice := slices[i]
		if !(len(slice) == length && reNumber.Match([]byte(slice))) {
			url := strings.Join(slices[:i+1], "/") + "/" + id
			return ctx.Redirect(http.StatusSeeOther, url)
		}
	}
	return errors.New("invalid url")
}

func (w *MultiStepForm) getSecurity(ctx echo.Context) GuidMaker {
	if w.Security != nil {
		return w.Security
	}
	return security.New()
}

func (w *MultiStepForm) getStorage(ctx echo.Context) echo.Storage {
	if w.Storage != nil {
		return w.Storage
	}
	return ctx.Echo().Cache
}

func makeMultiStepKey(id string) string {
	return "multistep:" + id
}

var (
	reNumber = regexp.MustCompile(`\d+`)
)

// Attach handler to MUX
func AddMultiStepHandler(
	mux Mux,
	handler echo.HandlerFunc,
) {
	// Undo multistep form
	mux.GET("/:id/:undo", handler)

	// Stage multistep form
	mux.Any("/:id", handler)

	// Start multistep form
	mux.GET("/", handler)
}

var redoDataSet = echo.NewDataSet(
	map[string]string{
		"redo": "redo",
	},
	false,
)
