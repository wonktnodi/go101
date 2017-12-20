package structs

import "fmt"

type Animal struct {
    Name string
}

type Human struct {
    *Animal
}

type AdvancedHuman struct {
    *Animal
}

type Pet struct{
    Name string
}

type Dog struct {
    *Animal
    *Pet
}

func (a *Animal) Speak() {
    fmt.Println("Animal speaks")
}

func (h *AdvancedHuman) Speak() {
    fmt.Println("Human speaks")
}

func Demo() {
    a := &Animal{
        Name: "Dennis",
    }

    b := &Human{a}
    c := &AdvancedHuman{a}
    d := &Dog{&Animal{Name: "Dennis"}, &Pet{Name: "Zedo"}}
    a.Speak() //Animal speaks
    b.Speak() //Animal speaks
    c.Speak() //Human speaks

    // name conflict
    fmt.Println(d.Animal.Name)
    fmt.Println(d.Pet.Name)
}
