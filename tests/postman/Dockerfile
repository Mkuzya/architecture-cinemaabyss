FROM node:18-alpine

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy test files
COPY . .

# Expose port for reports
EXPOSE 3000

# Default command
CMD ["node", "run-tests.js", "--environment", "docker", "--timeout", "30000"]