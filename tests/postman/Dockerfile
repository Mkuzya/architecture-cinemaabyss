FROM node:18-alpine

WORKDIR /app

# Install Newman
RUN npm install -g newman

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy test files
COPY . .

# Make script executable
RUN chmod +x run-tests.sh

# Default command
CMD ["npm", "run", "test:docker"]