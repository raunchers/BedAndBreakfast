package models

// Will never import other packages, exist to be imported itself

// TemplateData holds data sent from the handlers to the templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	// interface when you do not know what type the data may be
	Data      map[string]interface{}
	CSRFToken string
	//Flash message (like a success message)
	Flash   string
	Warning string
	Error   string
}
