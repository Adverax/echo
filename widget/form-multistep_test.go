package widget

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/adverax/echo"
	"github.com/adverax/echo/generic"
)

type stageBank struct {
	MultiStepBaseStage
}

func (stage *stageBank) Model(
	ctx echo.Context,
	state *MultiStepState,
) (echo.Model, error) {
	name := &FormTextInput{
		FormField: FormField{
			Name: "name",
		},
	}

	return echo.Model{
		"Name": name,
	}, nil
}

func (stage *stageBank) Consume(
	ctx echo.Context,
	state *MultiStepState,
	data interface{},
) (reply interface{}, err error) {
	return "amount", nil
}

type stageAmount struct {
	MultiStepBaseStage
	result generic.Params
}

func (stage *stageAmount) Model(
	ctx echo.Context,
	state *MultiStepState,
) (echo.Model, error) {
	amount := &FormTextInput{
		FormField: FormField{
			Name:  "amount",
			Codec: echo.UnsignedCodec,
		},
	}

	return echo.Model{
		"Amount": amount,
	}, nil
}

func (stage *stageAmount) Consume(
	ctx echo.Context,
	state *MultiStepState,
	data interface{},
) (reply interface{}, err error) {
	stage.result = state.Params
	return nil, ctx.Redirect(http.StatusSeeOther, "/index")
}

type payStrategy struct {
	bank   *stageBank
	amount *stageAmount
}

func (s *payStrategy) Setup(ctx echo.Context, values generic.Params) (stage string, err error) {
	return "bank", nil
}

func (s *payStrategy) Stage(ctx echo.Context, state *MultiStepState) (MultiStepStage, error) {
	if state.Stage == "amount" {
		return s.amount, nil
	}
	return s.bank, nil
}

func payHandler(strategy *payStrategy) echo.HandlerFunc {
	return echo.HandlerFunc(func(ctx echo.Context) error {
		form := &MultiStepForm{
			Form: Form{
				Model: make(echo.Model, 16),
			},
			Strategy: strategy,
			Timeout:  time.Hour,
		}

		state, err := form.Execute(ctx, nil)
		if err != nil || state == nil {
			if err == echo.ErrAbort {
				return nil
			}
			return err
		}

		res, err := form.Render(ctx)
		if err != nil {
			return err
		}

		layout := ` 
{{- with .Form -}}
<form name="payment">
{{- with .Model -}}
{{with .Name}}<input type="text" name="{{.Name}}" value="{{.Value}}">{{end}}
{{with .Amount}}<input type="text" name="{{.Name}}" value="{{.Value}}">{{end}}
{{if .Prev}}PrevBtn{{end -}}
{{if .Next}}NextBtn{{end -}}
{{- end}}
</end>
{{end -}}`
		layout += fmt.Sprintf("(-%s-)", state.Stage)
		tpl := template.Must(template.New("main").Parse(layout))
		return tpl.Execute(
			ctx.Response().Writer,
			map[string]interface{}{
				"Form": res,
			},
		)
	})
}

func TestMultiStep(t *testing.T) {
	e := echo.New()
	strategy := &payStrategy{
		amount: &stageAmount{
			MultiStepBaseStage: MultiStepBaseStage{
				PrevBtn: MultiStepPrevBtn,
			},
		},
		bank: &stageBank{
			MultiStepBaseStage: MultiStepBaseStage{
				NextBtn: MultiStepNextBtn,
			},
		},
	}
	handler := payHandler(strategy)

	// Init multiStep
	req := httptest.NewRequest(http.MethodGet, "/pay", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	err := handler(ctx)
	assert.NoError(t, err)
	require.Equal(t, http.StatusSeeOther, rec.Code)
	path := rec.HeaderMap["Location"][0]

	// Bank stage
	req = httptest.NewRequest(http.MethodGet, path, nil)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues(getParamMultiStepId(path))
	err = handler(ctx)
	assert.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "(-bank-)")
	require.Contains(t, rec.Body.String(), "NextBtn")
	require.NotContains(t, rec.Body.String(), "PrevBtn")

	// Submit bank page
	f := make(url.Values)
	f.Set("name", "UltraBank")
	req = httptest.NewRequest(http.MethodPost, path, strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues(getParamMultiStepId(path))
	err = handler(ctx)
	require.NoError(t, err)
	require.Equal(t, http.StatusSeeOther, rec.Code)

	// Amount stage
	req = httptest.NewRequest(http.MethodGet, path, nil)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues(getParamMultiStepId(path))
	err = handler(ctx)
	assert.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "(-amount-)")
	require.NotContains(t, rec.Body.String(), "NextBtn")
	require.Contains(t, rec.Body.String(), "PrevBtn")

	// Undo amount stage
	req = httptest.NewRequest(http.MethodGet, path, nil)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.SetParamNames("id", "undo")
	ctx.SetParamValues(getParamMultiStepId(path), "undo")
	err = handler(ctx)
	assert.NoError(t, err)
	require.Equal(t, http.StatusSeeOther, rec.Code)

	// Retry bank page
	req = httptest.NewRequest(http.MethodGet, path, nil)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues(getParamMultiStepId(path))
	err = handler(ctx)
	assert.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "(-bank-)")
	require.Contains(t, rec.Body.String(), "NextBtn")
	require.NotContains(t, rec.Body.String(), "PrevBtn")

	// Submit bank page
	f = make(url.Values)
	f.Set("name", "UltraBank2")
	req = httptest.NewRequest(http.MethodPost, path, strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues(getParamMultiStepId(path))
	err = handler(ctx)
	require.NoError(t, err)
	require.Equal(t, http.StatusSeeOther, rec.Code)

	// Retry amount stage
	req = httptest.NewRequest(http.MethodGet, path, nil)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues(getParamMultiStepId(path))
	err = handler(ctx)
	assert.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "(-amount-)")
	require.NotContains(t, rec.Body.String(), "NextBtn")
	require.Contains(t, rec.Body.String(), "PrevBtn")

	// Submit amount page
	f = make(url.Values)
	f.Set("amount", "25")
	req = httptest.NewRequest(http.MethodPost, path, strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues(getParamMultiStepId(path))
	err = handler(ctx)
	require.NoError(t, err)
	require.Equal(t, http.StatusSeeOther, rec.Code)
	assert.Equal(t,
		generic.Params{
			"name":   "UltraBank2",
			"amount": uint64(25),
		},
		strategy.amount.result,
	)
}

func getParamMultiStepId(s string) string {
	items := strings.Split(s, "/")
	return items[len(items)-1]
}
