@echo off
echo Menyalakan Database MySQL (Docker)...
docker-compose up -d

echo Menyalakan Backend Go...
start cmd /k "cd backend && go run main.go"

echo Menyalakan Frontend React...
start cmd /k "cd frontend && npm run dev"

echo Semua server sedang berjalan!
pause