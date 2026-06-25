# Student API

This is my first backend Go project. I made it to learn Go while building a simple Student API that performs CRUD operations on student info like name, email, and age. (NO AI used Btw)

I also built a simple React frontend to visualize and interact with the data!

## Project Structure

- **Main entry point:** `cmd/student-api/main.go`
- **Source code:** `internal/`
- **Database:** `storage/`
- **Config files:** `config/`

## How to Use

### Backend Setup

1. Clone the repository:
   ```bash
   git clone github.com/rusilkoirala/student-api
   ```

2. Create the database file at `/storage`:
   ```bash
   touch storage/storage.db
   ```

3. Run the backend from the project root:
   ```bash
   go run cmd/student-api/main.go -config config/local.yaml
   ```

   The backend should now be running.

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   # or
   pnpm install
   ```

3. Start the development server:
   ```bash
   npm run dev
   # or
   pnpm dev
   ```

4. Open your browser and go to `https://localhost:5173`

### Docker Deployment

You can run the backend in a container with a production-friendly setup:

```bash
docker build -t student-api .
docker run --rm -p 3000:3000 \
  -e CONFIG_PATH=/app/config/local.yaml \
  -v "$PWD/storage:/app/storage" \
  student-api
```

---

Thank you!