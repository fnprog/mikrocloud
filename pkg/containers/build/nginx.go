package build

import (
	"fmt"
)

func GenerateNginxConfig(isSPA bool) string {
	if isSPA {
		return `server {
    listen 80;
    server_name _;
    root /usr/share/nginx/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/x-javascript application/xml+rss application/json;
}`
	}

	return `server {
    listen 80;
    server_name _;
    root /usr/share/nginx/html;
    index index.html index.htm;

    location / {
        try_files $uri $uri/ =404;
    }

    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/x-javascript application/xml+rss application/json;
}`
}

func GenerateStaticDockerfile(config *StaticConfig) string {
	outputDir := config.OutputDir
	if outputDir == "" {
		outputDir = "dist"
	}

	var dockerfile string

	if config.BuildCommand != "" {
		dockerfile = fmt.Sprintf(`FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN %s

FROM nginx:alpine
COPY --from=builder /app/%s /usr/share/nginx/html
`, config.BuildCommand, outputDir)
	} else {
		dockerfile = fmt.Sprintf(`FROM nginx:alpine
COPY %s /usr/share/nginx/html
`, outputDir)
	}

	if config.NginxConfig != "" {
		dockerfile += fmt.Sprintf("COPY %s /etc/nginx/conf.d/default.conf\n", config.NginxConfig)
	} else {
		nginxConf := GenerateNginxConfig(config.IsSPA)
		dockerfile += fmt.Sprintf(`RUN echo '%s' > /etc/nginx/conf.d/default.conf
`, nginxConf)
	}

	dockerfile += `EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
`

	return dockerfile
}
