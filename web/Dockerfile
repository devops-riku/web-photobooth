# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files first to leverage Docker cache
COPY package*.json ./

# Install ALL requirements (dependencies + devDependencies) for building
RUN npm install

# Copy source code
COPY . .

# Build the SvelteKit application
RUN npm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

# Copy built artifacts
COPY --from=builder /app/build ./build
COPY --from=builder /app/package*.json ./

# Install ONLY production requirements to keep image small
RUN npm ci --omit=dev

# Environment setup
ENV NODE_ENV=production
ENV PORT=5173

EXPOSE 5173

# Start the Node.js server
CMD ["node", "build/index.js"]
