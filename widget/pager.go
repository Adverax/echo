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
	"net/url"
	"strconv"

	"github.com/adverax/echo"
	"github.com/adverax/echo/data"
)

type PagerInfo struct {
	CurPage      int // Current page number (the numbering starts from 1)
	LastPage     int // Last page number (the numbering starts from 1)
	FirstVisible int // First visible record number (the numbering starts from 1)
	LastVisible  int // Last visible record number (the numbering starts from 1)
	Total        int // Total records count
	Count        int // Visible records count
}

type PagerReport struct {
	Id      string                 `json:"id,omitempty"`
	Info    PagerInfo              `json:"info"`
	Buttons map[string]interface{} `json:"buttons,omitempty"`
	Message interface{}            `json:"message,omitempty"`
}

func (w *PagerReport) render(
	ctx echo.Context,
) interface{} {
	res := make(map[string]interface{}, 8)
	res["CurPage"] = w.Info.CurPage
	res["LastPage"] = w.Info.LastPage
	res["FirstVisible"] = w.Info.FirstVisible
	res["LastVisible"] = w.Info.LastVisible
	res["Total"] = w.Info.Total
	res["Count"] = w.Info.Count
	if len(w.Buttons) != 0 {
		res["Buttons"] = w.Buttons
	}
	res["Message"] = w.Message

	return res
}

type Pager struct {
	Id       string        // Pager identifier (optional)
	MsgEmpty Widget        // Message for empty record set (optional)
	MsgStats Widget        // Message for statistics )optional)
	Label    Widget        // Pager label (optional)
	Prev     Widget        // Previous page label (optional)
	Next     Widget        // Next page label (optional)
	Param    string        // Hot parameter (default "pg")
	Capacity int           // Items count per page (default 10). Without data provider it is row count
	BtnCount int           // Links count in pager (default 10)
	Url      *url.URL      // Base url (default used current request url)
	Provider data.Provider // Data provider
}

func (w *Pager) execute(
	ctx echo.Context,
) (*PagerReport, error) {
	capacity := w.Capacity
	if capacity == 0 {
		capacity = 10
	}

	if w.Provider == nil {
		message, err := Coalesce(ctx, w.MsgEmpty, MessageListNoRecords)
		if err != nil {
			return nil, err
		}

		return &PagerReport{
			Info: PagerInfo{
				CurPage:      1,
				LastPage:     1,
				FirstVisible: 1,
				LastVisible:  capacity,
				Total:        capacity,
				Count:        capacity,
			},
			Message: message,
		}, nil
	}

	btnCount := w.BtnCount
	if btnCount == 0 {
		btnCount = 10
	}

	param := w.Param
	if param == "" {
		param = "pg"
	}

	uu := w.Url
	if uu == nil {
		uu = new(url.URL)
		*uu = *ctx.Request().URL
	}

	prev := w.Prev
	if prev == nil {
		prev = MessagePagerPrev
	}

	next := w.Next
	if next == nil {
		next = MessagePagerNext
	}

	// Make pager info
	total, err := w.Provider.Total(ctx)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		message, err := Coalesce(ctx, w.MsgEmpty, MessageListNoRecords)
		if err != nil {
			return nil, err
		}

		return &PagerReport{
			Info: PagerInfo{
				CurPage:      1,
				LastPage:     1,
				FirstVisible: 0,
				LastVisible:  0,
				Total:        0,
				Count:        0,
			},
			Message: message,
		}, nil
	}

	var info PagerInfo
	params := uu.Query()
	{
		s := params.Get(param)
		pg, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			pg = 1
		}
		curPage := int(pg)
		if curPage < 1 {
			curPage = 1
		}
		lastPage := (total + capacity - 1) / capacity
		if curPage > lastPage {
			curPage = lastPage
		}

		firstVisible := (curPage-1)*capacity + 1
		lastVisible := curPage * capacity
		if lastVisible > total {
			lastVisible = total
		}

		info = PagerInfo{
			CurPage:      curPage,
			LastPage:     lastPage,
			FirstVisible: firstVisible,
			LastVisible:  lastVisible,
			Total:        total,
			Count:        lastVisible - firstVisible + 1,
		}

		err = w.Provider.Import(
			ctx,
			&data.Pagination{
				Offset: int64(capacity) * int64(info.CurPage-1),
				Limit:  int64(capacity),
			},
		)
		if err != nil {
			return nil, err
		}
	}

	message, err := Coalesce(ctx, w.MsgStats, MessageListRecords)
	if err != nil {
		return nil, err
	}

	msg, err := FormatMessage(ctx, message, info.FirstVisible, info.LastVisible, info.Total)
	if err != nil {
		return nil, err
	}

	if total <= capacity {
		return &PagerReport{
			Info:    info,
			Message: msg,
		}, nil
	}

	page := info.CurPage
	pageCount := (total + capacity - 1) / capacity
	start := page
	finish := page
	for rest := btnCount - 1; rest > 0; {
		var has bool
		if finish < pageCount {
			finish++
			rest--
			has = true
		}
		if start > 1 && rest != 0 {
			start--
			rest--
			has = true
		}
		if !has {
			break
		}
	}

	buttons := make(map[string]interface{}, 16)

	// Prev button
	label, err := prev.Render(ctx)
	if err != nil {
		return nil, err
	}
	params.Set(param, strconv.FormatInt(int64(page-1), 10))
	uu.RawQuery = params.Encode()
	action, err := RenderLink(ctx, uu)
	if err != nil {
		return nil, err
	}
	buttons["Prev"] = w.newBtn(label, action, page == 1)

	// Range of buttons
	if start <= finish {
		band := make([]map[string]interface{}, 0, finish-start+1)
		for pg := start; pg <= finish; pg++ {
			params.Set(param, strconv.FormatInt(int64(pg), 10))
			uu.RawQuery = params.Encode()
			action, err := RenderLink(ctx, uu)
			if err != nil {
				return nil, err
			}
			btn := make(map[string]interface{}, 4)
			btn["Label"] = strconv.FormatInt(int64(pg), 10)
			btn["Action"] = action
			if pg == page {
				btn["Active"] = true
			}
			band = append(band, btn)
		}
		buttons["Band"] = band
	}

	// Next button
	label, err = next.Render(ctx)
	if err != nil {
		return nil, err
	}
	params.Set(param, strconv.FormatInt(int64(page+1), 10))
	uu.RawQuery = params.Encode()
	action, err = RenderLink(ctx, uu)
	if err != nil {
		return nil, err
	}
	buttons["Next"] = w.newBtn(label, action, page >= pageCount)
	if len(buttons) == 0 {
		buttons = nil
	}

	return &PagerReport{
		Id:      w.Id,
		Info:    info,
		Buttons: buttons,
		Message: msg,
	}, nil
}

func (w *Pager) newBtn(
	label interface{},
	action interface{},
	disabled bool,
) map[string]interface{} {
	btn := make(map[string]interface{}, 4)
	btn["Label"] = label
	btn["Action"] = action
	if disabled {
		btn["Disabled"] = true
	}
	return btn
}
