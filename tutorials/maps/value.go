package maps

import "fmt"

type item struct {
    ID string
}

func MapValueMutable() {

    var data = map[int]*item{}
    v1 := item{"1"}
    v2 := item{"2"}
    v3 := item{"3"}
    data[1] = &v1
    data[2] = &v2
    data[3] = &v3

    v, ok := data[2]
    if ok {
        v.ID = "5"
        fmt.Printf("%p: %#v\n", v, v)
    }

    for k, v := range data {
        fmt.Printf("%v[%p]: %#v\n", k, v, v)
    }
}
