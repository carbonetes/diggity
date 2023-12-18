package cran

import (
	"strings"
)

type DescriptionFile struct {
	Package     string       `json:"package,omitempty"`
	Version     string       `json:"version,omitempty"`
	Date        string       `json:"date,omitempty"`
	Priority    string       `json:"priority,omitempty"`
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	Maintainers []Maintainer `json:"maintainers,omitempty"`
	Authors     []Author     `json:"authors,omitempty"`
	Depends     string       `json:"depends,omitempty"`
}

type Author struct {
	Name    string `json:"name,omitempty"`
	Role    string `json:"role,omitempty"`
	Email   string `json:"email,omitempty"`
	Comment string `json:"comment,omitempty"`
}

type Maintainer struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

func readManifestFile(content []byte) DescriptionFile {
	metadata := make(map[string]interface{})
	lines := strings.Split(string(content), "\n")
	var prev string
	for _, line := range lines {
		if strings.Contains(line, ": ") {
			parts := strings.Split(line, ": ")
			key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			metadata[key] = value
			prev = key
		} else {
			metadata[prev] = metadata[prev].(string) + " " + line
		}
	}
	if metadata["Author@"] != nil {
		delete(metadata, "Author@")
	}

	if metadata["Authors@R"] != nil {
		delete(metadata, "Authors@R")
	}
	var description DescriptionFile
	for key, value := range metadata {
		switch key {
		case "Package":
			description.Package = value.(string)
		case "Version":
			description.Version = value.(string)
		case "Date":
			description.Date = value.(string)
		case "Priority":
			description.Priority = value.(string)
		case "Title":
			description.Title = value.(string)
		case "Description":
			description.Description = value.(string)
		case "Maintainer":
			props := strings.Fields(value.(string))
			var name, email string
			for _, prop := range props {
				if !strings.Contains(prop, "@") {
					name += prop + " "
				} else {
					email = strings.Trim(prop, "<>")
				}
			}
			description.Maintainers = append(description.Maintainers, Maintainer{
				Name:  name,
				Email: email,
			})
		case "Author":
			props := strings.Split(value.(string), "\n")
			var authors []Author
			for _, prop := range props {
				parts := strings.Fields(prop)
				var name string
				var roles []string
				for _, part := range parts {
					if !strings.Contains(part, "[") {
						name += part + " "
					} else {
						role := strings.Replace(part, "[", "", -1)
						role = strings.Replace(role, "]", "", -1)
						role = strings.TrimSpace(role)
						if strings.Contains(role, ",") {
							r := strings.Split(role, ",")
							roles = append(roles, r...)
						} else {
							roles = append(roles, role)
						}
					}
				}
				authors = append(authors, Author{
					Name: strings.TrimSpace(name),
					Role: strings.Join(roles, ","),
				})
			}
			description.Authors = authors
		case "Depends":
			description.Depends = value.(string)
		}
	}

	return description
}
