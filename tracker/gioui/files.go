//go:build !js
// +build !js

package gioui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/vsariola/sointu"
)

func (t *Tracker) OpenSongFile(forced bool) {
	if !forced && t.ChangedSinceSave() {
		t.ConfirmSongActionType = ConfirmLoad
		t.ConfirmSongDialog.Visible = true
		return
	}
	if p := t.FilePath(); p != "" {
		d, _ := filepath.Split(p)
		d = filepath.Clean(d)
		t.OpenSongDialog.Directory.SetText(d)
		t.OpenSongDialog.FileName.SetText("")
	}
	t.OpenSongDialog.Visible = true
}

func (t *Tracker) SaveSongFile() bool {
	if p := t.FilePath(); p != "" {
		return t.saveSong(p)
	}
	t.SaveSongAsFile()
	return false
}

func (t *Tracker) SaveSongAsFile() {
	t.SaveSongDialog.Visible = true
	if p := t.FilePath(); p != "" {
		d, f := filepath.Split(p)
		d = filepath.Clean(d)
		t.SaveSongDialog.Directory.SetText(d)
		t.SaveSongDialog.FileName.SetText(f)
	}
}

func (t *Tracker) ExportWav() {
	t.ExportWavDialog.Visible = true
	if p := t.FilePath(); p != "" {
		d, _ := filepath.Split(p)
		d = filepath.Clean(d)
		t.ExportWavDialog.Directory.SetText(d)
	}
}

func (t *Tracker) LoadInstrument() {
	t.OpenInstrumentDialog.Visible = true
}

func (t *Tracker) SaveInstrument() {
	t.SaveInstrumentDialog.Visible = true
}

func (t *Tracker) loadSong(filename string) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	var song sointu.Song
	if errJSON := json.Unmarshal(b, &song); errJSON != nil {
		if errYaml := yaml.Unmarshal(b, &song); errYaml != nil {
			var err4kp error
			var patch sointu.Patch
			if patch, err4kp = sointu.Read4klangPatch(bytes.NewReader(b)); err4kp != nil {
				t.Alert.Update(fmt.Sprintf("Error unmarshaling a song file: %v / %v / %v", errYaml, errJSON, err4kp), Error, time.Second*3)
				return
			} else {
				song = t.Song()
				song.Score = t.Song().Score.Copy()
				song.Patch = patch
			}
		}
	}
	if song.Score.Length <= 0 || len(song.Score.Tracks) == 0 || len(song.Patch) == 0 {
		t.Alert.Update("The song file is malformed", Error, time.Second*3)
		return
	}
	t.SetSong(song)
	t.SetFilePath(filename)
	t.ClearUndoHistory()
	t.SetChangedSinceSave(false)
}

func (t *Tracker) saveSong(filename string) bool {
	var extension = filepath.Ext(filename)
	var contents []byte
	var err error
	if extension == ".json" {
		contents, err = json.Marshal(t.Song())
	} else {
		contents, err = yaml.Marshal(t.Song())
	}
	if err != nil {
		t.Alert.Update(fmt.Sprintf("Error marshaling a song file: %v", err), Error, time.Second*3)
		return false
	}
	if extension == "" {
		filename = filename + ".yml"
	}
	ioutil.WriteFile(filename, contents, 0644)
	t.SetFilePath(filename)
	t.SetChangedSinceSave(false)
	return true
}

func (t *Tracker) exportWav(filename string, pcm16 bool) {
	var extension = filepath.Ext(filename)
	if extension == "" {
		filename = filename + ".wav"
	}
	data, err := sointu.Play(t.synthService, t.Song(), true) // render the song to calculate its length
	if err != nil {
		t.Alert.Update(fmt.Sprintf("Error rendering the song during export: %v", err), Error, time.Second*3)
		return
	}
	buffer, err := sointu.Wav(data, pcm16)
	if err != nil {
		t.Alert.Update(fmt.Sprintf("Error converting to .wav: %v", err), Error, time.Second*3)
		return
	}
	ioutil.WriteFile(filename, buffer, 0644)
}

func (t *Tracker) saveInstrument(filename string) bool {
	var extension = filepath.Ext(filename)
	var contents []byte
	var err error
	if extension == ".json" {
		contents, err = json.Marshal(t.Instrument())
	} else {
		contents, err = yaml.Marshal(t.Instrument())
	}
	if err != nil {
		t.Alert.Update(fmt.Sprintf("Error marshaling a ínstrument file: %v", err), Error, time.Second*3)
		return false
	}
	if extension == "" {
		filename = filename + ".yml"
	}
	ioutil.WriteFile(filename, contents, 0644)
	return true
}

func (t *Tracker) loadInstrument(filename string) bool {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return false
	}
	var instrument sointu.Instrument
	if errJSON := json.Unmarshal(b, &instrument); errJSON != nil {
		if errYaml := yaml.Unmarshal(b, &instrument); errYaml != nil {
			var err4ki error
			if instrument, err4ki = sointu.Read4klangInstrument(bytes.NewReader(b)); err4ki != nil {
				t.Alert.Update(fmt.Sprintf("Error unmarshaling an instrument file: %v / %v / %v", errYaml, errJSON, err4ki), Error, time.Second*3)
				return false
			}
		}
	}
	// the 4klang instrument names are junk, replace them with the filename without extension
	instrument.Name = filepath.Base(filename[:len(filename)-len(filepath.Ext(filename))])
	if len(instrument.Units) == 0 {
		t.Alert.Update("The instrument file is malformed", Error, time.Second*3)
		return false
	}
	t.SetInstrument(instrument)
	if t.Instrument().Comment != "" {
		t.InstrumentEditor.ExpandComment()
	}
	return true
}
