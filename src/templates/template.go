package templates

import (
	"fmt"
	"net/http"

	"github.com/greatdanton/analytics/src/global"
)

// Execute is used to display template with templateName filled with data
// returns: error so we can return immediately in outer function
func Execute(w http.ResponseWriter, templateName string, data interface{}) error {
	err := global.Templates.ExecuteTemplate(w, templateName, data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return fmt.Errorf("Execute: ExecuteTemplate: %v", err)
	}
	return nil
}
