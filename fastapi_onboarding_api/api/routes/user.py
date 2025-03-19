from fastapi import APIRouter, HTTPException
from sqlmodel import select
from schemas.db_schemas import Users
from db_controller.db_controller import Session

user = APIRouter(prefix="/user", tags=["user_api"])

@user.get("/", response_model=list[Users], description="Get all users")
async def get_user(offset: int = 0, limit: int = 100):
    async with Session() as session:
        user = session.exec(select(Users).where(Users.role == "client").offset(offset).limit(limit)).all()
        if not user:
            raise HTTPException(status_code=404, detail={"error":{"404":"User is not found"}})
        return user