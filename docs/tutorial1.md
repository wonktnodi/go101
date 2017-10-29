#
很多编程语言都有字符串之间转换的机制，而GO语言则是通过模板来将一个对象的内容来作为参数传递从而字符串的转换。此方式不仅可以在重写HTML时插入对象值，也适用于其他方面。注意，本章内容并没有明确给出网络的工作方式，但对于网络编程方式很有用处。

### Introduction
### 介绍

大多数服务器端语言的机制主要是在静态页面插入一个动态生成的组件，如清单列表项目。典型的例子是在JSP、PHP和许多其他语言的脚本中。GO的template包中采取了相对简单的脚本化语言。
因为新的template包是刚刚被采用的，所有现在的template包中的文档少的可怜，旧的old/template包中也还存有少量的旧模板。新发布的帮助页面还没有关于新包的文档。关于template包的更改请参阅r60 (released 2011/09/07).
在这里，我们描述了这个新包。该包是描述了通过使用对象值改变了原来文本的方式从而在输入和输出时获取不同的文本。与JSP或类似的不同，它的作用不仅限于HTML文件，但在那可能会有更大的作用。
源文件被称作 template ，包括文本传输方式不变，以嵌入命令可以作用于和更改文本。命令规定如 {{ ... }} ，类似于JSP命令 <%= ... =%> 和PHP命令 <?php ... ?>。
Inserting object values

### 插入对象值

模板应用于GO对象中.GO对象的字段被插入到模板后，你就能从域中“挖”到他的子域，等等。当前对象以'.'代替, 所以把当前对象当做字符串插入时，你可以采用{{.}}的方式。这个包默认采用 fmt 包来作为插入值的字符串输出。

要插入当前对象的一个字段的值，你使用的字段名前加前缀 '.'。 例如, 如果要插入的对象的类型为


```
type Person struct {
        Name      string
        Age       int
        Emails     []string
        Jobs       []*Jobs
}

```

那么你要插入的字段 Name 和 Age 如下

```
The name is {{.Name}}.
The age is {{.Age}}.
```


我们可以使用range命令来循环一个数组或者链表中的元素。所以要获取 Emails 数组的信息，我们可以这么干


```
{{range .Emails}}
        ...
{{end}}

```

如果Job定义为

```
type Job struct {
    Employer string
    Role     string
}

```

如果我们想访问 Person字段中的 Jobs, 我们可以这么干 {{range .Jobs}}。这是一种可以将当前对象转化为Jobs 字段的方式. 通过 {{with ...}} ... {{end}} 这种方式, 那么{{.}} 就可以拿到Jobs 字段了,如下:

```
{{with .Jobs}}
    {{range .}}
        An employer is {{.Employer}}
        and the role is {{.Role}}
    {{end}}
{{end}}

```

你可以用这种方法操作任何类型的字段，而不仅限于数组。亲，用模板吧！

当我们拥有了模板,我们将它应用在对象中生成一个字符串，用这个对象来填充这个模板的值。分两步来实现模块的转化和应用，并且输出一个Writer, 如下

```
t := template.New("Person template")
t, err := t.Parse(templ)
if err == nil {
        buff := bytes.NewBufferString("")
        t.Execute(buff, person)
}

```

下面是一个例子来演示模块应用在对象上并且标准输入：

```
/**
 * PrintPerson
 */

package main

import (
        "fmt"
        "html/template"
        "os"
)

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

const templ = `The name is {{.Name}}.
The age is {{.Age}}.
{{range .Emails}}
        An email is {{.}}
{{end}}

{{with .Jobs}}
    {{range .}}
        An employer is {{.Employer}}
        and the role is {{.Role}}
    {{end}}
{{end}}
`

func main() {
        job1 := Job{Employer: "Monash", Role: "Honorary"}
        job2 := Job{Employer: "Box Hill", Role: "Head of HE"}

        person := Person{
                Name:   "jan",
                Age:    50,
                Emails: []string{"jan@newmarch.name", "jan.newmarch@gmail.com"},
                Jobs:   []*Job{&job1, &job2},
        }

        t := template.New("Person template")
        t, err := t.Parse(templ)
        checkError(err)

        err = t.Execute(os.Stdout, person)
        checkError(err)
}

func checkError(err error) {
        if err != nil {
                fmt.Println("Fatal error ", err.Error())
                os.Exit(1)
        }
}

```
输出如下：

```
The name is jan.
The age is 50.

        An email is jan@newmarch.name

        An email is jan.newmarch@gmail.com




        An employer is Monash
        and the role is Honorary

        An employer is Box Hill
        and the role is Head of HE


```

注意，上面有很多空白的输出，这是因为我们的模板中有很多空白。如果想消除它, 模板设置如下：

```
{{range .Emails}} An email is {{.}} {{end}}
```
在这个示例例中，我们用字符串应用于模板。你同样也可以用方法template.ParseFiles()来从文件中下载模板。因为某些原因，我还不没搞清楚(在早期版本没有强制要求),关联模板的名字必须要与文件列表的第一个文件的基名相同。话说，这个是BUG吗?

### Pipelines

管道


上述转换到模板中插入的文本块。这些字符基本上是任意的，是任何字符串的字段值。如果我们希望它们出现的是HTML文档（或其他的特殊形式）的一部分，那么我们将不得不脱离特定的字符序列。例如，要显示任意文本在HTML文档中，我们要将“<”改成“&lt”。GO模板有一些内建函数，其中之一是html。这些函数的作用与Unix的管道类似，从标准输入读取和写入到标准输出。

如果想用“.”来获取当前对象值并且应用于HTML转义，你可以在模板里写个“管道”:

```
{{. | html}}
```

其他方法类似。

Mike Samuel指出，目前在exp/template/html 包里有一个方便的方法。如果所有的模板中的条目需要通过html 模板函数，那么Go语言方法 Escape(t *template.Template)就能获取模板而后将html 函数添加到模板中不存在该函数的每个节点中。用于HTML文档的模板是非常有用的，并能在其他使用场合生成相似的方法模式。

### Defining functions
### 定义方法

模板使用对象化的字符串表示形式插入值，使用fmt包将对象转换为字符串。有时候，这并不是必需。例如，为了避免被垃圾邮件发送者掌握电子邮件地址，常见的方式是把字符号“@”替换为“at”，如“jan at newmarch.name”。如果我们要使用一个模板，显示在该表单中的电子邮件地址，那么我们就必须建立一个自定义的功能做这种转变。

每个模板函数中使用的模板本身有的一个名称，以及相关联的函数。他们用下面方式进行关联如下

```
type FuncMap map[string]interface{}
```

例如，如果我们希望我们的模板函数是“emailExpand”，用来关联到Go函数EmailExpander，然后，我们像这样添加函数到到模板中

```
t = t.Funcs(template.FuncMap{"emailExpand": EmailExpander})
```

EmailExpander通常像这样标记：

```
func EmailExpander(args ...interface{}) string
```


我们感兴趣的是在使用过程中，那是一个只有一个参数的函数，并且是个字符串。在Go模板库的现有功能有初步的代码来处理不符合要求的情况，所以我们只需要复制。然后，它就能通过简单的字符串操作来改变格式的电子邮件地址。程序如

```
/**
 * PrintEmails
 */

package main

import (
        "fmt"
        "os"
        "strings"
        "text/template"
)

type Person struct {
        Name   string
        Emails []string
}

const templ = `The name is {{.Name}}.
{{range .Emails}}
        An email is "{{. | emailExpand}}"
{{end}}
`

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

func main() {
        person := Person{
                Name:   "jan",
                Emails: []string{"jan@newmarch.name", "jan.newmarch@gmail.com"},
        }

        t := template.New("Person template")

        // add our function
 t = t.Funcs(template.FuncMap{"emailExpand": EmailExpander})

        t, err := t.Parse(templ)

        checkError(err)

        err = t.Execute(os.Stdout, person)
        checkError(err)
}

func checkError(err error) {
        if err != nil {
                fmt.Println("Fatal error ", err.Error())
                os.Exit(1)
        }
}

```
The output is
输出为：

```
The name is jan.

        An email is "jan at newmarch.name"

        An email is "jan.newmarch at gmail.com"
```

### Variables
### 变量


template包，允许您定义和使用变量。这样做的动机，可能我们会考虑通过把他们的名字当做电子邮件地址前缀打印出来。我们又使用这个类型

```
type Person struct {
        Name      string
        Emails     []string
}
```

为了访问email的所有字符串, 可以用 range，如下

```
{{range .Emails}}
    {{.}}
{{end}}
```

但是需要指出的是，我们无法用'.' 的形式来访问字段 Name，因为当他被转化成数组元素时，字段Name并不包括其中。解决方法是，将字段Name 存储为一个变量，那么它就能在任意范围内被访问。变量在模板中用法是加前缀'$'。所以可以这样

```
{{$name := .Name}}
{{range .Emails}}
    Name is {{$name}}, email is {{.}}
{{end}}
```


程序如下：

```
/**
 * PrintNameEmails
 */

package main

import (
        "html/template"
        "os"
        "fmt"
)

type Person struct {
        Name   string
        Emails []string
}

const templ = `{{$name := .Name}}
{{range .Emails}}
    Name is {{$name}}, email is {{.}}
{{end}}
`

func main() {
        person := Person{
                Name:   "jan",
                Emails: []string{"jan@newmarch.name", "jan.newmarch@gmail.com"},
        }

        t := template.New("Person template")
        t, err := t.Parse(templ)
        checkError(err)

        err = t.Execute(os.Stdout, person)
        checkError(err)
}

func checkError(err error) {
        if err != nil {
                fmt.Println("Fatal error ", err.Error())
                os.Exit(1)
        }
}
```

输出为

```
    Name is jan, email is jan@newmarch.name

    Name is jan, email is jan.newmarch@gmail.com
```

### Conditional statements
### 条件语句


继续我们那个Person的例子，假设我们只是想打印出来的邮件列表，而不关心其中的字段。我们可以用模板这么干

```
Name is {{.Name}}
Emails are {{.Emails}}
```
输出为：

```
Name is jan
Emails are [jan@newmarch.name jan.newmarch@gmail.com]
```

因为这个fmt包会显示一个列表。

在许多情况下，这样做也没有问题，如果那是你想要的。让我们考虑下一种情况，它 几乎是对的但不是必须的。有一个JSON序列化对象的包，让我们看看第4章。它是这样的

```
{"Name": "jan",
 "Emails": ["jan@newmarch.name", "jan.newmarch@gmail.com"]
}
```
JSON包是一个你会在实践中使用，但是让我们看看我们是否能够使用JSON输出模板。我们可以做一些我们有的类似的模板。这几乎就是一个JSON串行器：

```
{"Name": "{{.Name}}",
 "Emails": {{.Emails}}
}
```

像这样组装

```
{"Name": "jan",
 "Emails": [jan@newmarch.name jan.newmarch@gmail.com]
}
```


其中有两个问题：地址没有在引号中，列表中的元素应该是'，'分隔。

这样如何：在数组中的元素，把它们放在引号中并用逗号分隔？

```
{"Name": {{.Name}},
  "Emails": [
   {{range .Emails}}
      "{{.}}",
   {{end}}
  ]
}
```

It will produce
像这样组装

```
{"Name": "jan",
 "Emails": ["jan@newmarch.name", "jan.newmarch@gmail.com",]
}

```

(再加上一些空白)。

同样，这样貌似几乎是正确的，但如果你仔细看，你会看到尾有“，”在最后的列表元素。根据JSON的语法（请参阅 http://www.json.org/，这个结尾的'，'是不允许的。这样实现结果可能会有所不同。

我们想要打印所有在后面带','的元素除了最后一个。"这个确实有点难搞, 一个好方法"在',' 之前打印所有元素除了第一个。" (我在 "brianb"的 Stack Overflow上提了建议)。这样更易于实现，因为第一个元素索引为0，很多编程语言包括GO模板都将0当做布尔型的false。

条件语句的一种形式是{{if pipeline}} T1 {{else}} T0 {{end}}。我们需要通过pipeline来获取电子邮件到数组的索引。幸运的是， range的变化语句为我们提供了这一点。有两种形式，引进变量

```
{{range $elmt := array}}
{{range $index, $elmt := array}}
```


所以我们遍历数组，如果该索引是false（0），我们只是打印的这个索引的元素，否则打印它前面是','的元素。模板是这样的

```
{"Name": "{{.Name}}",
 "Emails": [
 {{range $index, $elmt := .Emails}}
    {{if $index}}
        , "{{$elmt}}"
    {{else}}
         "{{$elmt}}"
    {{end}}
 {{end}}
 ]
}
```


完整的程序如下

```
/**
 * PrintJSONEmails
 */

package main

import (
        "html/template"
        "os"
        "fmt"
)

type Person struct {
        Name   string
        Emails []string
}

const templ = `{"Name": "{{.Name}}",
 "Emails": [
{{range $index, $elmt := .Emails}}
    {{if $index}}
        , "{{$elmt}}"
    {{else}}
         "{{$elmt}}"
    {{end}}
{{end}}
 ]
}
`

func main() {
        person := Person{
                Name:   "jan",
                Emails: []string{"jan@newmarch.name", "jan.newmarch@gmail.com"},
        }

        t := template.New("Person template")
        t, err := t.Parse(templ)
        checkError(err)

        err = t.Execute(os.Stdout, person)
        checkError(err)
}

func checkError(err error) {
        if err != nil {
                fmt.Println("Fatal error ", err.Error())
                os.Exit(1)
        }
}
```

上面给出的是正确的JSON输出

在结束本节之前，我们强调了用逗号分隔的列表格式的问题，解决方式是可以在模板函数中定义适当的函数。正如俗话说的，“道路不止一条！”下面的程序是Roger Peppe给我的：

```
/**
 * Sequence.go
 * Copyright Roger Peppe
 */

package main

import (
        "errors"
        "fmt"
        "os"
        "text/template"
)

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

func main() {
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
```
### Conclusion
### 结论


template包在对于某些类型的文本转换涉及插入对象值的情况是非常有用的。虽然它没有正则表达式功能强大，但它执行比正则表达式速度更快，在许多情况下比正则表达式更容易使用。

