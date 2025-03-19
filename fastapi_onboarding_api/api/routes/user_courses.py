from db_controller.db_controller import Session
from fastapi import APIRouter, HTTPException
from sqlmodel import  select
from schemas.db_schemas import UserCourses, UserCoursesBase
from schemas.schemas_update import UserCoursesUpdate

user_courses = APIRouter(prefix="/user/courses", tags=["user_courses_api"])

@user_courses.post("/create", description="Create a new course for user", response_model=UserCoursesBase)
async def create_user_courses(user_course: UserCoursesBase):
    async with Session() as session:
        db_course = UserCourses.model_validate(user_course)
        session.add(db_course)
        session.commit()
        session.refresh(db_course)
        return db_course

@user_courses.get("/get",description="Get all courses and users")
async def get_user_courses(offset: int = 0, limit: int = 100):
    async with Session() as session:
        user_courses = session.exec(select(UserCourses).offset(offset).limit(limit)).all()
        if not user_courses:
            raise HTTPException(status_code=404, detail={"error":{"404":"Courses and users not created yet"}})
        return user_courses

@user_courses.get("/get/{course_id}", description="Get a specific user & course by id")
async def get_user_course_by_id(course_id: int):
    async with Session() as session:
        course = session.get(UserCourses, course_id)
        if not course:
            raise HTTPException(status_code=404, detail={"error":{"404":"Course and user is not created yet"}})
        return course

@user_courses.delete("/delete/{course_id}", description="Delete a course")
async def delete_user_courses(course_id: int):
    async with Session() as session:
        course = session.get(UserCourses, course_id)
        if not course:
            raise HTTPException(status_code=404, detail={"error":{"404":"Course and user is not created yet"}})
        session.delete(course)
        session.commit()
        return "User course deleted"

@user_courses.patch("/update/{course_id}", description="Update a course & course by id")
async def update_user_courses(course_id: int, user_course: UserCoursesUpdate):
    async with Session() as session:
        course_db = session.get(UserCourses, course_id)
        if not course_db:
            raise HTTPException(status_code=404, detail={"error":{"404":"Course and user is not created yet"}})
        course_data = user_course.model_dump(exclude_unset=True)
        course_db.sqlmodel_update(course_data)
        session.add(course_db)
        session.commit()
        session.refresh(course_db)
        return course_db