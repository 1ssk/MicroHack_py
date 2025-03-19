from fastapi import APIRouter
from sqlmodel import SQLModel, create_engine
from sqlmodel import Session as DBSession
from contextlib import asynccontextmanager
from dotenv import load_dotenv
from sqlalchemy.exc import OperationalError
from time import sleep
from loguru import logger

import os

load_dotenv()

db_user = os.getenv("DB_USER")
db_pwd = os.getenv("DB_PASSWORD")
db_host = os.getenv("DB_HOST")
db_port = os.getenv("DB_PORT")
db_name = os.getenv("DB_NAME")

DATABASE_URL = f'postgresql://{db_user}:{db_pwd}@{db_host}:{db_port}/{db_name}'
engine = create_engine(DATABASE_URL, pool_reset_on_return=None, pool_recycle=3600)


#Создаем сессию для БД которая будет нормально открываться и закрываться.
@asynccontextmanager
async def Session():
    session = DBSession(engine)
    try:
        logger.info("Starting session")
        yield session
        logger.info("Session finished")
    except Exception as e:
        logger.error(f"Error {e}")
        sleep(5)
    finally:
        session.close()


@asynccontextmanager
async def lifespan(app: APIRouter):
    SQLModel.metadata.create_all(engine)
    yield
    engine.dispose()

db_control = APIRouter(prefix="/db_controller", tags=["db_controller"], lifespan=lifespan)