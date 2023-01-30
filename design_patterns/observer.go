package main

import "fmt"

type Topic interface {
	register(observer Observer)
	broadcast()
}
type Observer interface {
	getId() string
	updataValue(string)
}

// Item  -> no avaible
// Item -> alert -> hay items
type Item struct {
	observer  []Observer
	name      string
	aveilable bool
}

func NewItem(name string) *Item {
	return &Item{
		name: name,
	}
}

func (i *Item) UpdateAvaible() {
	fmt.Printf("Item %s is available \n", i.name)
	i.aveilable = true
	i.broadCast()
}

func (i *Item) broadCast() {
	for _, observer := range i.observer {
		observer.updataValue(i.name)
	}
}

func (i *Item) register(observer Observer) {
	i.observer = append(i.observer, observer)
}

type EmailCLient struct {
	id string
}

func (eC EmailCLient) getId() string {
	return eC.id
}

func (eC EmailCLient) updataValue(value string) {
	fmt.Printf("Sending email - %s available from client %s\n", value, eC.id)
}
func main() {
	nvidiaItem := NewItem("rtx")
	firstObserver := &EmailCLient{
		id: "12ab",
	}
	secondObserver := &EmailCLient{
		id: "34dc",
	}
	nvidiaItem.register(firstObserver)
	nvidiaItem.register(secondObserver)
	nvidiaItem.UpdateAvaible()
}
