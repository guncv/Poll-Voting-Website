# Stage 1: Build the app
FROM node:23-slim AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build
RUN npm prune --production

# Stage 2: Production image
FROM node:23-slim

WORKDIR /app

COPY --from=builder /app ./

EXPOSE 3000

ENV HOST=0.0.0.0
ENV PORT=3000
CMD ["npm", "start"]
