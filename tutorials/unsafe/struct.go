package unsafe

import (
    "unsafe"
    "fmt"
)

type Person struct {
    name   string
    age    int
    gender bool
}

type Combine struct {
    Person
    country string
    city    string
}

func UnsafeDemo() {
    a := [4]int{0, 1, 2, 3}
    p1 := unsafe.Pointer(&a[1])
    p3 := unsafe.Pointer(uintptr(p1) + 2*unsafe.Sizeof(a[0]))
    *(*int)(p3) = 6
    fmt.Println("a =", a)

    who := Person{"John", 30, true}
    pp := unsafe.Pointer(&who)
    pname := (*string)(unsafe.Pointer(uintptr(pp) + unsafe.Offsetof(who.name)))
    page := (*int)(unsafe.Pointer(uintptr(pp) + unsafe.Offsetof(who.age)))
    pgender := (*bool)(unsafe.Pointer(uintptr(pp) + unsafe.Offsetof(who.gender)))
    *pname = "Alice"
    *page = 28
    *pgender = false
    fmt.Println(who)

    //illegalUseA()
    //illegalUseB()
    structConvert()
}

func structConvert() {
    var combine Combine
    combine.age = 30
    combine.name = "Alice"
    combine.gender = false
    combine.city = "shanghai"
    combine.country = "china"

    basePtr := (*Person)(unsafe.Pointer(&combine))
    fmt.Println(basePtr)

    base := *(*Person)(unsafe.Pointer(&combine))
    fmt.Println(base)

}

// case A: conversions between unsafe.Pointer and uintptr
//         don't appear in the same expression
func illegalUseA() {
    fmt.Println("===================== illegalUseA")

    pa := new([4]int)

    // split the legal use
    // p1 := unsafe.Pointer(uintptr(unsafe.Pointer(pa)) + unsafe.Sizeof(pa[0]))
    // into two statements (illegal use):
    ptr := uintptr(unsafe.Pointer(pa))
    p1 := unsafe.Pointer(ptr + unsafe.Sizeof(pa[0]))
    // "go vet" will make a warning for the above line:
    // possible misuse of unsafe.Pointer

    // the unsafe package docs, https://golang.org/pkg/unsafe/#Pointer,
    // thinks above splitting is illegal.
    // but the current Go compiler and runtime (1.8) can't detect
    // this illegal use.
    //
    // however, to make your program run well absolutely,
    // it is best to comply with the unsafe package docs.

    *(*int)(p1) = 123
    fmt.Println("*(*int)(p1)  :", *(*int)(p1)) //
}

// case B: pointers are pointing at unknown addresses
func illegalUseB() {
    fmt.Println("===================== illegalUseB")

    a := [4]int{0, 1, 2, 3}
    p := unsafe.Pointer(&a)
    p = unsafe.Pointer(uintptr(p) + uintptr(len(a))*unsafe.Sizeof(a[0]))
    // now p is pointing at the end of the memory occupied by value a.
    // up to now, although p is illegal, it is no problem.
    // but it would be dangerous if we modify the value pointed by p
    *(*int)(p) = 123
    fmt.Println("*(*int)(p)  :", *(*int)(p)) // 123 or not 123
    // the current Go compiler/runtime (1.8) and "go vet"
    // will not detect the illegal use here.

    // however, the current Go runtime (1.8) will
    // detect the illegal use and panic for the below code.
    p = unsafe.Pointer(&a)
    for i := 0; i <= len(a); i++ {
        *(*int)(p) = 123 // Go runtime (1.8) never panic here in the tests

        fmt.Println(i, ":", *(*int)(p))
        // panic at the above line for the last iteration, when i==4.
        // runtime error: invalid memory address or nil pointer dereference

        p = unsafe.Pointer(uintptr(p) + unsafe.Sizeof(a[0]))
    }
}
