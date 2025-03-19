package handlers

import (
	"html/template"
	"log"
	"net/http"
	"sync"
)

// TemplateCache — глобальный кэш шаблонов
var TemplateCache = map[string]*template.Template{}
var cacheMutex sync.RWMutex

// InitTemplateCache загружает шаблоны в память
func InitTemplateCache(templateDir string, templates []string) {
	for _, tmpl := range templates {
		path := templateDir + "/" + tmpl
		t, err := template.ParseFiles(path)
		if err != nil {
			log.Fatalf("Ошибка загрузки шаблона %s: %v", tmpl, err)
		}
		cacheMutex.Lock()
		TemplateCache[tmpl] = t
		cacheMutex.Unlock()
		log.Printf("Шаблон %s успешно загружен", tmpl)
	}
}

// renderTemplate рендерит шаблон из кэша
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) error {
	cacheMutex.RLock()
	t, ok := TemplateCache[tmpl]
	cacheMutex.RUnlock()
	if !ok {
		http.Error(w, "Template not found", http.StatusNotFound)
		return nil
	}
	return t.Execute(w, data)
}

// handleError обрабатывает ошибки
func handleError(w http.ResponseWriter, err error, statusCode int) {
	log.Printf("Ошибка: %v", err)
	http.Error(w, http.StatusText(statusCode), statusCode)
}
