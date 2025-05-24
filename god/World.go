package god

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type Entity struct {
	Name       string
	Properties map[string]string
}

type World struct {
	entities map[string]Entity
	lEvents  []string
	mu       sync.RWMutex
}
type DumpResponse struct {
	Entities map[string]map[string]string // map of entity name → properties
	Events   []string
}
type GetEntityResponse struct {
	Entity   Entity `json:"entity"`
	ChaosMsg string `json:"chaos_msg,omitempty"`
}

func NewWorld() *World {
	w := &World{
		entities: make(map[string]Entity),
	}
	go w.mutateEntities()
	return w
}

func chaosChance(p float32) bool {
	return rand.Float32() < p
}

func (w *World) Create(name string, props map[string]string) (chaosMsg string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if chaosChance(0.3) {
		switch rand.Intn(3) {
		case 0:
			name = getRandomName()
			fmt.Printf("Chaos Created %s instead.\n", name)
			chaosMsg = fmt.Sprintf("Chaos Created %s instead.\n", name)
			w.lEvents = append(w.lEvents, chaosMsg)

		case 1:
			if rand.Float32() < 0.5 {
				randName := w.getRandomExistingName()
				delete(w.entities, randName)
				fmt.Printf("Chaos Deleted %s during creation.\n", randName)
				chaosMsg = fmt.Sprintf("Chaos Deleted %s during creation.\n", randName)
				w.lEvents = append(w.lEvents, chaosMsg)

			}
		case 2:
			randName := w.getRandomExistingName()
			if randName != "" {
				entity := w.entities[randName]
				entity.Properties["form"] = getRandomMutation()
				w.entities[randName] = entity
				fmt.Printf("Chaos Mutated %s during creation.\n", randName)
				chaosMsg = fmt.Sprintf("Chaos Mutated %s during creation.\n", randName)
				w.lEvents = append(w.lEvents, chaosMsg)

			}
		}
	}

	w.entities[name] = Entity{Name: name, Properties: props}
	w.lEvents = append(w.lEvents, fmt.Sprintf(" %s was created.\n", name))

	return chaosMsg
}

func (w *World) Get(name string) (Entity, bool, string) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var chaosMsg string

	if chaosChance(0.3) {
		name = w.getRandomExistingName()
		fmt.Printf("Chaos shuffled things around and you got %s instead.\n", name)
		chaosMsg = fmt.Sprintf("Chaos shuffled things around and you got %s instead.\n", name)
		w.lEvents = append(w.lEvents, chaosMsg)

	}
	e, ok := w.entities[name]
	return e, ok, chaosMsg
}

func (w *World) Destroy(name string) (chaosMsg string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if chaosChance(0.3) {
		switch rand.Intn(2) {
		case 0:
			name = w.getRandomExistingName()
			fmt.Printf("Chaos Destroyed %s instead.\n", name)
			chaosMsg = fmt.Sprintf("Chaos Destroyed %s instead.\n", name)
			w.lEvents = append(w.lEvents, chaosMsg)

		case 1:
			newName := getRandomName()
			props := map[string]string{"origin": "chaos"}
			w.entities[newName] = Entity{Name: newName, Properties: props}
			fmt.Printf("Chaos Created %s during destruction.\n", newName)
			chaosMsg = fmt.Sprintf("Chaos Created %s during destruction.\n", newName)
			w.lEvents = append(w.lEvents, chaosMsg)

		}
	}

	delete(w.entities, name)
	w.lEvents = append(w.lEvents, fmt.Sprintf("%s was destroyed", name))

	return chaosMsg
}

func (w *World) Dump() DumpResponse {
	w.mu.RLock()
	defer w.mu.RUnlock()

	entities := make(map[string]map[string]string)
	for name, entity := range w.entities {
		entities[name] = entity.Properties
	}

	// Assuming w.lEvents is []string
	eventsCopy := make([]string, len(w.lEvents))
	copy(eventsCopy, w.lEvents)

	return DumpResponse{
		Entities: entities,
		Events:   eventsCopy,
	}
}

func getRandomName() string {
	names := []string{"Tree", "Mountain", "River", "Fire", "Bird"}
	return names[rand.Intn(len(names))]
}

func getRandomMutation() string {
	forms := []string{"stone", "fire", "water", "fog", "light"}
	return forms[rand.Intn(len(forms))]
}

func (w *World) getRandomExistingName() string {
	for k := range w.entities {
		return k
	}
	return ""
}

func combineNameWithForm(name, form string) string {
	lowerName := strings.ToLower(name)
	if strings.Contains(lowerName, "corrupted") {
		return fmt.Sprintf("Even More Corrupted %s", name) // already corrupted gets more corrupted

	}
	if strings.Contains(lowerName, "of") || strings.Contains(lowerName, "on") {
		return fmt.Sprintf("Corrupted %s", name) // already mutated gets corrupted
	}
	if rand.Intn(2) == 0 {
		return fmt.Sprintf("%s of %s", name, form)
	}
	return fmt.Sprintf("%s on %s", name, form)
}

func (w *World) mutateEntities() {
	for {
		time.Sleep(10 * time.Second)
		w.mu.Lock()
		for name, entity := range w.entities {
			if chaosChance(0.3) {
				form := getRandomMutation()
				newName := combineNameWithForm(entity.Name, form)
				entity.Name = newName
				entity.Properties["mutated"] = "true"
				entity.Properties["form"] = form
				delete(w.entities, name)
				w.entities[newName] = entity
				fmt.Printf("Chaos mutated %s into %s.\n", name, newName)
				chaosMsg := fmt.Sprintf("Chaos mutated %s into %s.\n", name, newName)
				w.lEvents = append(w.lEvents, chaosMsg)

			}
		}
		w.mu.Unlock()
	}
}

func (w *World) AllEntities() map[string]Entity {
	w.mu.RLock()
	defer w.mu.RUnlock()

	copy := make(map[string]Entity)
	for k, v := range w.entities {
		copy[k] = v
	}
	return copy
}
