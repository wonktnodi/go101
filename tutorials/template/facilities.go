package template

import (
    "html/template"
    "log"
    "os"
    "fmt"
    "strings"
    "errors"
)

func add(left int, right int) int {
    return left + right
}

func check(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

type Person struct {
    Name   string
    Age    int
    Emails []string
    Jobs   []*Job
}

type Job struct {
    Employer string
    Role     string
}

func Simple() {
    tpl := `The name is {{.Name}}.
The age is {{.Age}}.
{{range .Emails}}An email is {{.}} {{end}}
{{with .Jobs}}{{range .}}An employer is {{.Employer}}
and the role is {{.Role}}{{end}}
{{end}}
`

    job1 := Job{Employer: "Monash", Role: "Honorary"}
    job2 := Job{Employer: "Box Hill", Role: "Head of HE"}

    person := Person{
        Name:   "jan",
        Age:    50,
        Emails: []string{"jan@newmarch.name", "jan.newmarch@gmail.com"},
        Jobs:   []*Job{&job1, &job2},
    }

    t, err := template.New("simple").Parse(tpl)
    check(err)
    err = t.Execute(os.Stdout, person)
    check(err)
}

func EmailExpander(args ...interface{}) string {
    ok := false
    var s string
    if len(args) == 1 {
        s, ok = args[0].(string)
    }
    if !ok {
        s = fmt.Sprint(args...)
    }

    // find the @ symbol
    substrs := strings.Split(s, "@")
    if len(substrs) != 2 {
        return s
    }
    // replace the @ by " at "
    return (substrs[0] + " at " + substrs[1])
}

func Pipeline() {
    const templ = `The name is {{.Name}}.
{{range .Emails}}An email is "{{. | emailExpand}}"
{{end}}
`
    person := Person{
        Name:   "jan",
        Emails: []string{"jan@newmarch.name", "jan.newmarch@gmail.com"},
    }

    t := template.New("Person template")
    // add our function
    t = t.Funcs(template.FuncMap{"emailExpand": EmailExpander})

    t, err := t.Parse(templ)
    check(err)
    err = t.Execute(os.Stdout, person)
    check(err)
}

func Variables() {
    const templ = `{{$name := .Name}}
{{range .Emails}}
    Name is {{$name}}, email is {{.}}
{{end}}
`
    person := Person{
        Name:   "jan",
        Emails: []string{"jan@newmarch.name", "jan.newmarch@gmail.com"},
    }

    t := template.New("Person template")
    t, err := t.Parse(templ)
    check(err)

    err = t.Execute(os.Stdout, person)
    check(err)
}

func Conditional() {
    const templ = `{"Name": "{{.Name}}",
 "Emails": [
   {{range $index, $elmt := .Emails}}{{if $index}},
   "{{$elmt}}"{{else}}"{{$elmt}}"{{end}}{{end}}
 ]
}
`
    person := Person{
        Name:   "jan",
        Emails: []string{"jan@newmarch.name", "jan.newmarch@gmail.com"},
    }

    t := template.New("Person template")
    t, err := t.Parse(templ)
    check(err)

    err = t.Execute(os.Stdout, person)
    check(err)
}

var tmpl = `{{$comma := sequence "" ", "}}
{{range $}}{{$comma.Next}}{{.}}{{end}}
{{$comma := sequence "" ", "}}
{{$colour := cycle "black" "white" "red"}}
{{range $}}{{$comma.Next}}{{.}} in {{$colour.Next}}{{end}}
`
var fmap = template.FuncMap{
    "sequence": sequenceFunc,
    "cycle":    cycleFunc,
}

type generator struct {
    ss []string
    i  int
    f  func(s []string, i int) string
}

func (seq *generator) Next() string {
    s := seq.f(seq.ss, seq.i)
    seq.i++
    return s
}

func sequenceGen(ss []string, i int) string {
    if i >= len(ss) {
        return ss[len(ss)-1]
    }
    return ss[i]
}

func cycleGen(ss []string, i int) string {
    return ss[i%len(ss)]
}

func sequenceFunc(ss ...string) (*generator, error) {
    if len(ss) == 0 {
        return nil, errors.New("sequence must have at least one element")
    }
    return &generator{ss, 0, sequenceGen}, nil
}

func cycleFunc(ss ...string) (*generator, error) {
    if len(ss) == 0 {
        return nil, errors.New("cycle must have at least one element")
    }
    return &generator{ss, 0, cycleGen}, nil
}

func Complex() {
    t, err := template.New("").Funcs(fmap).Parse(tmpl)
    if err != nil {
        fmt.Printf("parse error: %v\n", err)
        return
    }
    err = t.Execute(os.Stdout, []string{"a", "b", "c", "d", "e", "f"})
    if err != nil {
        fmt.Printf("exec error: %v\n", err)
    }
}