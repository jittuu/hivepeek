package peek

import (
  "io"
  "html/template"
)

const layout = "templates/layout.html"

func RenderTemplate(w io.Writer, data interface{}, filenames ...string) {
	filenames = append(filenames, layout)
	if t, err := template.New("layout.html").ParseFiles(filenames...); err != nil {
		panic(err)
	} else {
		if tErr := t.Execute(w, data); tErr != nil {
			panic(tErr)
		}
	}
}
