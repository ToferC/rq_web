package main

import "fmt"

type History struct {
	Name    string
	Event   *Event
	head    *Event
	current *Event
}

type Event struct {
	Name      string
	Year      int
	next      *Event
	Modifiers []*Modifier
}

type Modifier struct {
	Name string
}

func createHistory(name string) *History {
	return &History{
		Name: name,
	}
}

func (h *History) addEvent(name string, year int) error {
	e := &Event{
		Name: name,
		Year: year,
	}

	if h.head == nil {
		h.head = e
	} else {
		currentNode := h.head
		for currentNode.next != nil {
			currentNode = currentNode.next
		}
		currentNode.next = e
	}
	return nil
}

func (h *History) showAllEvents() error {
	currentNode := h.head
	if currentNode == nil {
		fmt.Println("History is empty")
		return nil
	}
	fmt.Printf("%+v\n", *currentNode)
	for currentNode.next != nil {
		currentNode = currentNode.next
		fmt.Printf("%+v\n", currentNode)
	}
	return nil
}

func (h *History) startHistory() *Event {
	h.current = h.head
	return h.current
}

func (h *History) nextEvent() *Event {
	h.current = h.current.next
	return h.current
}

func (h *History) showCurrent() {
	fmt.Printf("Current Event: %s - %d\n", h.current.Name, h.current.Year)
}

func main() {
	history := createHistory("Sartar 1583")
	history.addEvent("Civil War", 1584)
	history.addEvent("Uprising", 1585)
	history.addEvent("Starvation", 1586)
	history.addEvent("Flight", 1587)
	history.addEvent("Rebuilding", 1588)
	history.addEvent("Partnerships", 1589)
	history.addEvent("New Futures", 1590)

	//history.showAllEvents()

	fmt.Println()

	history.startHistory()
	history.showCurrent()

	history.nextEvent()
	history.showCurrent()

	history.nextEvent()
	history.showCurrent()

	history.nextEvent()
	history.showCurrent()

	history.nextEvent()
	history.showCurrent()

}
