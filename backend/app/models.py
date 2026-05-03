from datetime import datetime

from sqlalchemy import Column, DateTime, ForeignKey, Integer, String, Table
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.database import Base

api_key_secrets = Table(
    "api_key_secrets",
    Base.metadata,
    Column("api_key_id", Integer, ForeignKey("api_keys.id", ondelete="CASCADE"), primary_key=True),
    Column("secret_id", Integer, ForeignKey("secrets.id", ondelete="CASCADE"), primary_key=True),
)


class Secret(Base):
    __tablename__ = "secrets"

    id: Mapped[int] = mapped_column(Integer, primary_key=True, autoincrement=True)
    key: Mapped[str] = mapped_column(String, unique=True, nullable=False)
    value: Mapped[str] = mapped_column(String, nullable=False)
    created_at: Mapped[datetime] = mapped_column(DateTime, default=datetime.now, nullable=False)

    api_keys: Mapped[list["ApiKey"]] = relationship(
        secondary=api_key_secrets,
        back_populates="secrets",
    )


class ApiKey(Base):
    __tablename__ = "api_keys"

    id: Mapped[int] = mapped_column(Integer, primary_key=True, autoincrement=True)
    key: Mapped[str] = mapped_column(String, unique=True, nullable=False)
    created_at: Mapped[datetime] = mapped_column(DateTime, default=datetime.now, nullable=False)

    secrets: Mapped[list["Secret"]] = relationship(
        secondary=api_key_secrets,
        back_populates="api_keys",
    )
