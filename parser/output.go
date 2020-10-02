/*
 * Copyright 2019 The CovenantSQL Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package parser

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type reportCookieRecord struct {
	Name         string
	Path         string
	Domain       string
	Expires      time.Time
	MaxAge       int
	Expiry       string
	Secure       bool
	HttpOnly     bool
	SameSite     string
	UsedRequests int
	Category     string
	Description  string

	URL        string
	RemoteAddr string
	Status     int
	MimeType   string
	Initiator  string
	Source     string
	LineNo     int
}

type reportRecord struct {
	Category    string
	Description string
	Cookies     []*reportCookieRecord
}

type reportData struct {
	ScanTime        time.Time
	ScanURL         string
	CookieCount     int
	ScreenShotImage string
	Records         []*reportRecord
}

func (t *Task) OutputJSON(pretty bool) (str string, err error) {
	var jsonBlob []byte
	if pretty {
		jsonBlob, err = json.MarshalIndent(t.reportData, "", "  ")
	} else {
		jsonBlob, err = json.Marshal(t.reportData)
	}
	str = string(jsonBlob)
	return
}

func (t *Task) OutputHTML() (str string, err error) {
	return outputAsHTML(t.reportData)
}

func (t *Task) OutputPDF() (blob []byte, err error) {
	var f *os.File
	if f, err = ioutil.TempFile("", "gdpr_cookie*.html"); err != nil {
		return
	}

	tempHTML := f.Name()
	defer func() {
		_ = os.Remove(tempHTML)
	}()

	htmlData, err := outputAsHTML(t.reportData)
	if err != nil {
		return
	}

	_, _ = f.WriteString(htmlData)
	_ = f.Sync()
	_ = f.Close()

	return outputAsPDF(t.remote, tempHTML)
}

func (t *Task) OutputPDFToFile(filename string) (err error) {
	bt, err := t.OutputPDF()
	if err != nil {
		return
	}

	err = ioutil.WriteFile(filename, bt, 0644)

	return
}

func (t *Task) FormatEmail() (str string, err error) {
	return formatEmailContent(t.reportData)
}
