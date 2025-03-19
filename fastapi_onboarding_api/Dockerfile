FROM python:3.13-alpine AS builder

COPY --from=ghcr.io/astral-sh/uv:0.6.3 /uv /bin

WORKDIR /temp

RUN apk add --no-cache postgresql-dev gcc musl-dev postgresql-libs

ENV UV_COMPILE_BYTECODE=1 UV_LINK_MODE=copy

RUN --mount=type=cache,target=/root/.cache/uv \
    --mount=type=bind,source=uv.lock,target=uv.lock \
    --mount=type=bind,source=pyproject.toml,target=pyproject.toml \
    uv sync --frozen --no-install-project --no-dev


FROM python:3.13-alpine
WORKDIR /fastapi

COPY --from=builder /temp/.venv /fastapi/.venv

ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    PATH="/fastapi/.venv/bin:$PATH"

COPY . .

CMD ["python", "app.py"]

