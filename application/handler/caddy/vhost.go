/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package caddy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"strings"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/filemanager"
	"github.com/admpub/nging/application/library/modal"
	"github.com/admpub/nging/application/library/notice"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/mvc/events"
)

func init() {
	handler.RegisterToGroup(`/manage`, func(g *echo.Group) {
		g.Route(`GET`, ``, VhostIndex)
		g.Route(`GET,POST`, `/vhost_add`, VhostAdd)
		g.Route(`GET,POST`, `/vhost_edit`, VhostEdit)
		g.Route(`GET,POST`, `/vhost_delete`, VhostDelete)
		g.Route(`GET,POST`, `/vhost_file`, VhostFile)
		g.Route(`GET`, `/vhost_build`, Vhostbuild)
		g.Route(`GET`, `/clear_cache`, ClearCache)
	})
}

func VhostIndex(ctx echo.Context) error {
	m := model.NewVhost(ctx)
	page, size, totalRows, p := handler.PagingWithPagination(ctx)
	cnt, err := m.List(nil, nil, page, size)
	ret := handler.Err(ctx, err)
	if totalRows <= 0 {
		totalRows = int(cnt())
		p.SetRows(totalRows)
	}
	ctx.Set(`pagination`, p)
	ctx.Set(`listData`, m.Objects())
	return ctx.Render(`manage/index`, ret)
}

func Vhostbuild(ctx echo.Context) error {
	saveFile, err := getSaveDir()
	if err == nil {
		err = filepath.Walk(saveFile, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return filepath.SkipDir
			}
			return os.Remove(path)
		})
	}
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(`/manage`)
	}
	m := model.NewVhost(ctx)
	n := 100
	cnt, err := m.ListByOffset(nil, nil, 0, n, `disabled`, `N`)
	for i, j := 0, cnt(); int64(i) < j; i += n {
		if i > 0 {
			_, err = m.ListByOffset(nil, nil, i, n, `disabled`, `N`)
			if err != nil {
				handler.SendFail(ctx, err.Error())
				return ctx.Redirect(`/manage`)
			}
		}
		for _, m := range m.Objects() {
			var formData url.Values
			err := json.Unmarshal([]byte(m.Setting), &formData)
			if err == nil {
				file := filepath.Join(saveFile, fmt.Sprint(m.Id)+`.conf`)
				err = saveVhostConf(ctx, file, formData)
			}
			if err != nil {
				handler.SendFail(ctx, err.Error())
				return ctx.Redirect(`/manage`)
			}
		}
	}
	err = config.DefaultCLIConfig.CaddyRestart()
	if err != nil {
		ctx.Logger().Error(err)
	}
	handler.SendOk(ctx, ctx.T(`操作成功`))
	return ctx.Redirect(`/manage`)
}

func VhostAdd(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		m := model.NewVhost(ctx)
		m.Domain = ctx.Form(`domain`)
		m.Disabled = ctx.Form(`disabled`)
		m.Root = ctx.Form(`root`)
		var b []byte
		b, err = json.Marshal(ctx.Forms())
		switch {
		case err == nil:
			m.Setting = string(b)
			_, err = m.Add()
			if err != nil {
				break
			}
			fallthrough
		case 0 == 1:
			err = saveVhostData(ctx, m.Vhost, ctx.Forms())
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(`/manage`)
		}

		ctx.SetFunc(`Val`, func(name, defaultValue string) string {
			return ctx.Form(name)
		})
	} else {
		ctx.SetFunc(`Val`, func(name, defaultValue string) string {
			return defaultValue
		})
	}
	return ctx.Render(`manage/vhost_edit`, err)
}

func getSaveDir() (saveFile string, err error) {
	saveFile, err = filepath.Abs(config.DefaultConfig.Sys.VhostsfileDir)
	if err != nil {
		return
	}
	if fi, er := os.Stat(saveFile); er != nil || !fi.IsDir() {
		err = os.MkdirAll(saveFile, 0666)
		if err != nil {
			return
		}
	}
	return
}

func saveVhostConf(ctx echo.Context, saveFile string, values url.Values) error {
	SetCaddyfileFunc(ctx, values)
	ctx.Set(`values`, values)
	b, err := ctx.Fetch(`manage/caddyfile`, nil)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(saveFile, b, os.ModePerm)
	return err
}

func saveVhostData(ctx echo.Context, m *dbschema.Vhost, values url.Values) (err error) {
	var saveFile string
	saveFile, err = getSaveDir()
	if err != nil {
		return
	}
	saveFile = filepath.Join(saveFile, fmt.Sprint(m.Id)+`.conf`)
	if m.Disabled == `Y` {
		err = os.Remove(saveFile)
		if os.IsNotExist(err) {
			err = nil
		}
	} else {
		err = saveVhostConf(ctx, saveFile, values)
		if len(ctx.Form(`restart`)) > 0 {
			err = config.DefaultCLIConfig.CaddyRestart()
		}
	}
	return
}

func VhostDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	if id < 1 {
		handler.SendFail(ctx, ctx.T(`id无效`))
		return ctx.Redirect(`/manage`)
	}
	m := model.NewVhost(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err != nil {
		handler.SendFail(ctx, err.Error())
	} else {
		var saveFile string
		saveFile, err = filepath.Abs(config.DefaultConfig.Sys.VhostsfileDir)
		if err == nil {
			saveFile = filepath.Join(saveFile, fmt.Sprint(id)+`.conf`)
			err = os.Remove(saveFile)
			if os.IsNotExist(err) {
				err = nil
			}
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
			}
		}
	}
	return ctx.Redirect(`/manage`)
}

func VhostEdit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	if id < 1 {
		handler.SendFail(ctx, ctx.T(`id无效`))
		return ctx.Redirect(`/manage`)
	}

	var err error
	m := model.NewVhost(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(`/manage`)
	}
	if ctx.IsPost() {
		m.Domain = ctx.Form(`domain`)
		m.Disabled = ctx.Form(`disabled`)
		m.Root = ctx.Form(`root`)
		var b []byte
		b, err = json.Marshal(ctx.Forms())
		switch {
		case err == nil:
			m.Setting = string(b)
			err = m.Edit(nil, db.Cond{`id`: id})
			if err != nil {
				break
			}
			fallthrough
		case 0 == 1:
			err = saveVhostData(ctx, m.Vhost, ctx.Forms())
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(`/manage`)
		}
	} else {
		var formData url.Values
		if e := json.Unmarshal([]byte(m.Setting), &formData); e == nil {
			for key, values := range formData {
				for _, v := range values {
					ctx.Request().Form().Add(key, v)
				}
			}
		}
	}
	ctx.SetFunc(`Val`, func(name, defaultValue string) string {
		return ctx.Form(name)
	})
	ctx.Set(`activeURL`, `/manage`)
	return ctx.Render(`manage/vhost_edit`, err)
}

func ClearCache(ctx echo.Context) error {
	if err := modal.Clear(); err != nil {
		return err
	}
	notice.Clear()
	events.Event(`clearCache`, func(_ bool) {})
	return ctx.String(ctx.T(`已经清理完毕`))
}

func VhostFile(ctx echo.Context) error {
	var err error
	vhostId := ctx.Formx(`id`).Uint()
	filePath := ctx.Form(`path`)
	do := ctx.Form(`do`)
	m := model.NewVhost(ctx)
	err = m.Get(nil, db.Cond{`id`: vhostId})
	mgr := filemanager.New(m.Root, config.DefaultConfig.Sys.EditableFileMaxBytes, ctx)
	absPath := m.Root
	if err == nil && len(m.Root) > 0 {
		var exit bool

		if len(filePath) > 0 {
			filePath = filepath.Clean(filePath)
			absPath = filepath.Join(m.Root, filePath)
		}

		switch do {
		case `edit`:
			data := ctx.Data()
			if _, ok := Editable(absPath); !ok {
				data.SetInfo(ctx.T(`此文件不能在线编辑`), 0)
			} else {
				content := ctx.Form(`content`)
				encoding := ctx.Form(`encoding`)
				dat, err := mgr.Edit(absPath, content, encoding)
				if err != nil {
					data.SetInfo(err.Error(), 0)
				} else {
					data.SetData(dat, 1)
				}
			}
			return ctx.JSON(data)
		case `rename`:
			data := ctx.Data()
			newName := ctx.Form(`name`)
			err = mgr.Rename(absPath, newName)
			if err != nil {
				data.SetInfo(err.Error(), 0)
			} else {
				data.SetCode(1)
			}
			return ctx.JSON(data)
		case `delete`:
			err = mgr.Remove(absPath)
			if err != nil {
				handler.SendFail(ctx, err.Error())
			}
			return ctx.Redirect(ctx.Referer())
		case `upload`:
			err = mgr.Upload(absPath)
			if err != nil {
				user := handler.User(ctx)
				if user != nil {
					notice.OpenMessage(user.Username, `upload`)
					notice.Send(user.Username, notice.NewMessageWithValue(`upload`, ctx.T(`文件上传出错`), err.Error()))
				}
				return err
			}
			return ctx.String(`OK`)
		default:
			var dirs []os.FileInfo
			err, exit, dirs = mgr.List(absPath)
			ctx.Set(`dirs`, dirs)
		}
		if exit {
			return err
		}
	}
	ctx.Set(`data`, m)
	if filePath == `.` {
		filePath = ``
	}
	ctx.Set(`path`, filePath)
	ctx.Set(`absPath`, absPath)
	ctx.SetFunc(`Editable`, func(fileName string) bool {
		_, ok := Editable(fileName)
		return ok
	})
	ctx.Set(`activeURL`, `/manage`)
	return ctx.Render(`manage/file`, err)
}

func Editable(fileName string) (string, bool) {
	ext := strings.TrimPrefix(filepath.Ext(fileName), `.`)
	ext = strings.ToLower(ext)
	typ, ok := config.DefaultConfig.Sys.EditableFileExtensions[ext]
	return typ, ok
}
