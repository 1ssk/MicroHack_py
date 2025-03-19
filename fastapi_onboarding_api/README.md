**Документация по запуску FastAPI**
1. Установите UV-виртуальное окружение следующей командой в ***Powershell*** Windows:
```powershell
powershell -ExecutionPolicy ByPass -c "irm https://astral.sh/uv/install.ps1 | iex"
```
Либо для ***Linux/MacOS***
```bash
  curl -LsSf https://astral.sh/uv/install.sh | sh
```
2. Создайте виртуальное окружение внутри папки с проектом (FastAPI) командой:
```bash
  uv venv .venv #Создание виртуального окружения 
```
3. Войдите в виртуальное окружение для Windows:
```shell
  .\.venv\Scripts\activate
```
Либо для ***Linux/MacOS***:
```bash
    source .venv/bin/activate
```
4. Синхронизируйте все зависимости проекта:
```bash
  uv sync 
```
5. Запустите FastAPI следующими вариантами:
```bash
    1. uvicorn app:app --reload --port 8000 
    2. python app.py
```
**Docker**
```shell
    docker compose build
    docker compose up
```
