from db_controller.db_controller import Session
from fastapi import APIRouter, HTTPException
from sqlmodel import  select
from schemas.db_schemas import Lessons, LessonsBase
from schemas.schemas_update import LessonsUpdate

lessons = APIRouter(prefix="/lessons", tags=["lessons_api"])

@lessons.post("/create", description="Create lesson", response_model=LessonsBase)
async def lessons_create(lesson: LessonsBase):
    async with Session() as session:
        db_lessons = Lessons.model_validate(lesson)
        session.add(db_lessons)
        session.commit()
        session.refresh(db_lessons)
        return db_lessons


@lessons.delete("/delete/{lessons_id}", description="Delete lesson by id")
async def lessons_delete(lessons_id: int):
    async with Session() as session:
        lesson = session.get(Lessons, lessons_id)
        if not lesson:
            raise HTTPException(status_code=404, detail="Lesson not found")
        session.delete(lesson)
        session.commit()
        return "Lesson deleted"


@lessons.get("/get/", description="Get all lessons with limit")
async def lessons_get(offset: int = 0, limit: int = 100):
    async with Session() as session:
        lesson =  session.exec(select(Lessons).offset(offset).limit(limit)).all()
        if not lesson:
            raise HTTPException(status_code=404, detail={"error":{"404":"Lesson is not found"}})
        return lesson


@lessons.get("/get/{lesson_id}/", description="Get a lesson by id")
async def lesson_get(lesson_id: int):
    async with Session() as session:
        lesson = session.get(Lessons, lesson_id)
        if not lesson:
            raise HTTPException(status_code=404, detail={"error":{"404":"Lesson is not found"}})
        return lesson


@lessons.patch("/update/{lesson_id}", description="Update a lesson by id")
async def lessons_update(lesson_id: int, lesson: LessonsUpdate):
    async with Session() as session:
        db_lesson = session.get(Lessons, lesson_id)
        if not db_lesson:
            raise HTTPException(status_code=404, detail={"error":{"404":"Lesson is not found"}})
        lesson_data = lesson.model_dump(exclude_unset=True)
        db_lesson.sqlmodel_update(lesson_data)
        session.add(db_lesson)
        session.commit()
        session.refresh(db_lesson)
        return db_lesson