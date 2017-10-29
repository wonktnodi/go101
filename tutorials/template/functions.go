package template

import (
    "strings"
    "html/template"
)

var (
    templateMap = template.FuncMap{
        "Upper": upperString,
    }
)

func upperString(s string) string {
    return strings.ToUpper(s)
}