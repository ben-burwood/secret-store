import base64
import secrets as pysecrets

from robyn import Request, Robyn
from sqlalchemy import delete, select

from app.database import get_session
from app.models import ApiKey, api_key_secrets


def _is_global():
    return ApiKey.id.not_in(select(api_key_secrets.c.api_key_id))


def register_api_key_routes(app: Robyn):
    @app.get("/web/api/key", auth_required=True)
    async def get_global_api_key(request: Request):
        with get_session() as session:
            row = session.scalars(select(ApiKey).where(_is_global()).order_by(ApiKey.created_at.desc())).first()
            return {"key": row.key if row else None}

    @app.get("/web/api/key/generate", auth_required=True)
    async def generate_global_api_key(request: Request):
        new_key = base64.b64encode(pysecrets.token_bytes(32)).decode("ascii")
        with get_session() as session:
            session.execute(delete(ApiKey).where(_is_global()))
            session.add(ApiKey(key=new_key))
            session.commit()
        return {"key": new_key}
