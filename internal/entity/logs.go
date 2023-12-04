package entity

import (
	"bytes"
	"sort"
	"strconv"
	"time"
)

type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type LabelSet []Label

type Log struct {
	FormatedLog string   `json:"formatedLog"`
	Labels      LabelSet `json:"labels"`
}

type LogFilters struct {
	ProductID    string
	VersionTag   string
	From         time.Time
	To           time.Time
	WorkflowName string
	ProcessName  string
	RequestID    string
	Level        string
	Logger       string
	Limit        int
}

type LogOutFormat string

const (
	OutFormatConsole LogOutFormat = "console"
	OutFormatFile    LogOutFormat = "file"
)

func (of LogOutFormat) IsValid() bool {
	return of == OutFormatConsole || of == OutFormatFile
}

func (ls LabelSet) String() string {
	var b bytes.Buffer

	labelSetMap := make(map[string]string)
	for _, label := range ls {
		labelSetMap[label.Key] = label.Value
	}

	keys := make([]string, 0, len(labelSetMap))
	for k := range labelSetMap {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	b.WriteByte('{')

	for i, k := range keys {
		if i > 0 {
			b.WriteByte(',')
			b.WriteByte(' ')
		}

		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(strconv.Quote(labelSetMap[k]))
	}

	b.WriteByte('}')

	return b.String()
}
