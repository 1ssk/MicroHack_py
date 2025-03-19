from sqlmodel import SQLModel, Field

class Users(SQLModel, table=True):
    __tablename__ = 'users'
    id: int | None = Field(default=None, primary_key=True, index=True)
    username: str = Field(nullable=False)
    password: str = Field(nullable=False)
    role: str = Field(nullable=False)


class LessonsBase(SQLModel):
    title: str = Field(nullable=True)
    course_id: int = Field(nullable=True, foreign_key="courses.id")
    url: str = Field(nullable=True)


class Lessons(LessonsBase, table=True):
    __tablename__ = 'lessons'
    id: int | None = Field(default=None, primary_key=True, index=True)


class UserCoursesBase(SQLModel):
    user_id: int = Field(default=None, foreign_key="users.id")
    course_id: int = Field(default=None, foreign_key="courses.id")


class UserCourses(UserCoursesBase, table=True):
    __tablename__ = 'user_courses'
    id: int | None = Field(default=None, primary_key=True, exclude=True)


class CoursesBase(SQLModel):
    title: str = Field(nullable=True)

class Courses(CoursesBase, table=True):
    __tablename__ = 'courses'
    id: int | None = Field(default=None, primary_key=True, index=True)
