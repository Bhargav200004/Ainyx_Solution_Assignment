# ----- Build Stage -----
FROM node:22-alpine AS builder

WORKDIR /app

# Cache dependencies.
COPY package.json package-lock.json ./
RUN npm ci

# Copy source and build.
COPY . .
RUN npm run build

# ----- Runtime Stage -----
FROM nginx:alpine

# Remove default nginx config.
RUN rm /etc/nginx/conf.d/default.conf

# Copy custom config and built assets.
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
