# Build stage
FROM node:22-alpine AS builder

WORKDIR /app

# Install build tools for native modules
RUN apk add --no-cache python3 make g++

COPY package*.json ./
RUN npm ci --legacy-peer-deps

COPY . .

# Point the webapp at the nginx-proxied API path
ENV VITE_API_URL=/api
RUN npm run build

# Serve stage
FROM nginx:alpine

COPY --from=builder /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
