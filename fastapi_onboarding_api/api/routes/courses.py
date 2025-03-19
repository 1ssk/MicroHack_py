from fastapi import APIRouter, HTTPException
from schemas.db_schemas import Courses, CoursesBase
from schemas.schemas_update import CoursesUpdate
from sqlmodel import select
from db_controller.db_controller import Session


courses = APIRouter(prefix="/courses", tags=["courses_api"])

@courses.post("/create", response_model=CoursesBase, description="Create a new course")
async def create_course(course: CoursesBase):
    async with Session() as session:
        db_course = Courses.model_validate(course)
        session.add(db_course)
        session.commit()
        session.refresh(db_course)
        return db_course


@courses.get("/get", description="Get list courses with limit")
async def get_courses(offset: int = 0, limit: int = 100):
    async with Session() as session:
        course = session.exec(select(Courses).offset(offset).limit(limit)).all()
        if not course:
            raise HTTPException(status_code=404, detail={"error":{"404":"Course is not found"}})
        return course


@courses.get("/get/{course_id}", description="Get course by id")
async def get_course_by_id(course_id: int):
    async with Session() as session:
        course = session.get(Courses, course_id)
        if not course:
            raise HTTPException(status_code=404, detail={"error":{"404":"Course is not found"}})
        return course


@courses.delete("/delete/{course_id}", description="Delete course by id")
async def delete_course_by_id(course_id: int):
    async with Session() as session:
        course = session.get(Courses, course_id)
        if not course:
            raise HTTPException(status_code=404, detail={"error":{"404":"Course is not found"}})
        session.delete(course)
        session.commit()
        return "Course deleted"


@courses.patch("/update/{course_id}", description="Update course by id")
async def update_course(course_id: int, course: CoursesUpdate):
    async with Session() as session:
        db_course = session.get(Courses, course_id)
        if not db_course:
            raise HTTPException(status_code=404, detail={"error":{"404":"Course is not found"}})
        course_data = course.model_dump(exclude_unset=True)
        db_course.sqlmodel_update(course_data)
        session.add(db_course)
        session.commit()
        session.refresh(db_course)
        return db_course