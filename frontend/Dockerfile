# Stage 1: Build the app
FROM node:18 AS builder

WORKDIR /app

# Copy only package files first to leverage Docker cache for dependency installation.
COPY package*.json ./
RUN npm install

# Now copy the rest of the source code.
COPY . .

# Build the production assets (adjust this if your build output is in a different folder)
RUN npm run build

# Optionally remove dev dependencies (if your build output doesn't need them)
RUN npm prune --production

# Stage 2: Production image
FROM node:18-alpine

WORKDIR /app

# Copy the production-ready files from the builder stage.
COPY --from=builder /app ./

EXPOSE 3000
CMD ["npm", "start"]
