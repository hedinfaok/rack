package manifest

import (
	"fmt"
	"strings"
)

func (m *Manifest) Build(dir string, s Stream, noCache bool) error {
	pulls := map[string]string{}
	builds := []Service{}

	for _, service := range m.Services {
		dockerFile := service.Build.Dockerfile
		if dockerFile == "" {
			dockerFile = service.Dockerfile
		}
		switch {
		case service.Build.Context != "":
			builds[fmt.Sprintf("%s|%s", service.Build.Context, coalesce(dockerFile, "Dockerfile"))] = service.Tag()
		case service.Image != "":
			pulls[service.Image] = service.Tag()
		}
	}

	for build, tag := range builds {
		parts := strings.SplitN(build, "|", 2)

		args := []string{"build"}

		args = append(args, "-f", parts[1])
		args = append(args, "-t", tag)
		args = append(args, parts[0])

		run(s, Docker(args...))
		// runPrefix(systemPrefix(m), Docker(args...))
	}

	for image, tag := range pulls {
		args := []string{"pull"}

		args = append(args, image)

		run(s, Docker(args...))
		run(s, Docker("tag", image, tag))
		// runPrefix(systemPrefix(m), Docker(args...))
		// runPrefix(systemPrefix(m), Docker("tag", image, tag))
	}

	return nil
}
