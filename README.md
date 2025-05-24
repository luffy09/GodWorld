# 🌍 GodWorld

**GodWorld** is a playful, chaos-driven world simulator API built in Go. It allows you to create, retrieve, mutate, and destroy entities — but beware! A random chance of "chaos" can cause unpredictable behavior, such as random creations, deletions, or mutations.

It also features a Swagger-powered UI for exploring the API interactively.

---

## 🚀 Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/your-username/GodWorld.git
cd GodWorld
```

### 2. To Start

```bash
go mod tidy
go run main.go
```

If swagger doesn't open by itself here is a link
http://localhost:8080/swagger/index.html

### 3. Example Create Input
```bash
{
  "name": "Planet",
  "properties": {
    "Surface": "Rock",
    "Atmosphere": "Hydrogen"
  }
}
```
