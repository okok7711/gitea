package bruhman

import (
        "code.gitea.io/gitea/modules/markup"
	"code.gitea.io/gitea/modules/markup/markdown"
)

// Render takes a markdown string and returns the rendered HTML string or an error.
func Render(markdownContent string) (string, error) {
	buffer, err := markdown.RenderString(&markup.RenderContext{
		Ctx: nil,
		Links: markup.Links{
			Base: "http://localhost:3000/user13/repo11/",
		},
	}, markdownContent)
	if err != nil {
		return "", err
	}

	return string(buffer), nil
}
