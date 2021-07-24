package xhtml

import "fmt"

func ModalTextField(title, extra, id string) string {
	if extra != "" {
		extra += " "
	}
	return fmt.Sprintf("<div class=\"row align-items-center\" style=\"margin: 8px;\">\n<div class=\"col-auto\"><span>%s %s</span></div>\n<div class=\"col\"><input type=\"text\" id=\"%s\"></div>\n</div>", title, extra, id)
}

func ModalNumberField(title, extra, id string, min, max, step int) string {
	if extra != "" {
		extra += " "
	}
	return fmt.Sprintf("<div class=\"row align-items-center\" style=\"margin: 8px;\">\n<div class=\"col-auto\"><span>%s %s</span></div>\n<div class=\"col\"><input type=\"number\" id=\"%s\" min=\"%d\" max=\"%d\" step=\"%d\"></div>\n</div>", title, extra, id, min, max, step)
}

func ModalSelectSingleField(title, extra, id string, values [][]string, selected int) string {
	if extra != "" {
		extra += " "
	}
	options := ""
	for i, v := range values {
		if i == selected {
			options += fmt.Sprintf("<option value=\"%s\" selected=\"\">%s</option>\n", v[1], v[0])
		} else {
			options += fmt.Sprintf("<option value=\"%s\">%s</option>\n", v[1], v[0])
		}
	}
	return fmt.Sprintf("<div class=\"row align-items-center\" style=\"margin: 8px;\">\n<div class=\"col-auto\"><span>%s %s</span></div>\n<div class=\"col\"><select id=\"%s\">\n%s</select></div>\n</div>", title, extra, id, options)
}
