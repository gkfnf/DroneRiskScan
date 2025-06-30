module.exports = {
  apps: [
    {
      name: "drone-scanner-frontend",
      script: "npm",
      args: "start",
      cwd: "./",
      instances: 1,
      autorestart: true,
      watch: false,
      max_memory_restart: "1G",
      env: {
        NODE_ENV: "production",
        PORT: 3000,
      },
    },
    {
      name: "drone-scanner-backend",
      script: "./backend/main",
      cwd: "./",
      instances: 1,
      autorestart: true,
      watch: false,
      max_memory_restart: "512M",
      env: {
        GO_ENV: "production",
        PORT: 8080,
      },
    },
  ],
}
