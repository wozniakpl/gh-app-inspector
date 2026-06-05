package inspector

import (
	"reflect"
	"sort"
	"strings"

	"github.com/google/go-github/v88/github"
)

func permissionRows(p *github.InstallationPermissions) [][2]string {
	if p == nil {
		return nil
	}
	v := reflect.ValueOf(p).Elem()
	t := v.Type()
	var rows [][2]string
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() != reflect.Ptr || f.IsNil() {
			continue
		}
		s, ok := f.Elem().Interface().(string)
		if !ok {
			continue
		}
		rows = append(rows, [2]string{jsonName(t.Field(i).Tag.Get("json")), s})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i][0] < rows[j][0] })
	return rows
}

func jsonName(tag string) string {
	if tag == "" {
		return ""
	}
	if i := strings.Index(tag, ","); i >= 0 {
		return tag[:i]
	}
	return tag
}
