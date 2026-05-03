from cryptography.fernet import Fernet
from sqlalchemy.types import String, TypeDecorator

from app.config import ENCRYPTION_KEY

try:
    _fernet = Fernet(ENCRYPTION_KEY.encode("ascii"))
except ValueError as e:
    raise RuntimeError(f"ENCRYPTION_KEY is invalid: {e}") from e


class EncryptedString(TypeDecorator):
    impl = String
    cache_ok = True

    def process_bind_param(self, value, dialect):
        if value is None:
            return None
        return _fernet.encrypt(value.encode("utf-8")).decode("ascii")

    def process_result_value(self, value, dialect):
        if value is None:
            return None
        return _fernet.decrypt(value.encode("ascii")).decode("utf-8")
