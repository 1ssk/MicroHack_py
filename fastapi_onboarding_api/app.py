from api.routes.lessons import lessons
from api.routes.user import user
from api.routes.courses import courses
from api.routes.user_courses import user_courses
from db_controller.db_controller import db_control
from fastapi import FastAPI
from fastapi.responses import JSONResponse
from fastapi.exceptions import ResponseValidationError
from fastapi.middleware.cors import CORSMiddleware
from sqlalchemy.exc import IntegrityError

import uvicorn

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"]
)


#Обработка исключений
@app.exception_handler(IntegrityError)
async def integrity_error_handler(request, exc):
    return JSONResponse(
        status_code=500,
        content={"error": {"500":"Foreign Key is not found, correct it!"}},
    )

@app.exception_handler(ResponseValidationError)
async def response_validation_error(request, exc):
    return JSONResponse(
        status_code=500,
        content={"error": {"500":"Input should be a valid dictionary or object to extract fields from"}},
    )


#Подключаем маршуты для user
app.include_router(user)
#Подключаем маршруты для courses
app.include_router(courses)
#Подключаем контроллер БД
app.include_router(db_control)
#Подключаем роут для добавления уроков
app.include_router(lessons)
#Подключаем роут для привзяки курсов к пользователям
app.include_router(user_courses)

if __name__ == "__main__":
    uvicorn.run("app:app", host="0.0.0.0", port=8000, log_level="info")