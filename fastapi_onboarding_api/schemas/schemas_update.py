from sqlmodel import SQLModel


class UsersUpdate(SQLModel):
    id: int | None = None
    username: str | None = None
    password: str | None = None
    role: str | None = None

class LessonsUpdate(SQLModel):
    title: str | None = None
    course_id: int | None = None
    url: str | None = None


class CoursesUpdate(SQLModel):
    title: str | None = None


class UserCoursesUpdate(SQLModel):
    id: int | None = None
    user_id: int | None = None
    course_id: int | None = None