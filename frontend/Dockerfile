# Stage 1: Build the app
FROM node:23-slim AS builder

WORKDIR /app

# Copy only package files first to leverage Docker cache for dependency installation.
COPY package*.json ./
RUN npm install

# Now copy the rest of the source code excluding files like node_modules and build artifacts.
COPY . .

# Build the production assets
RUN npm run build

# Optionally remove dev dependencies
RUN npm prune --production

# Stage 2: Production image
FROM node:23-slim

WORKDIR /app

# Copy the production-ready files from the builder stage.
COPY --from=builder /app ./

EXPOSE 3000
CMD ["npm", "start"]
