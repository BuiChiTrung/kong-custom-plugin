# Use an official Go runtime with Gin installed as the base image
FROM golang:latest

# Set the working directory
WORKDIR /fiber-app

# Copy the application source files to the container
COPY . .

# Expose port 8080 for the application to listen on
EXPOSE 9080

# Define the command to run the application
CMD ["go", "run", "."]
