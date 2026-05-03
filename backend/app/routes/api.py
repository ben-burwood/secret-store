from urllib.parse import unquote

from robyn import Request, Response, Robyn
from sqlalchemy import select
from sqlalchemy.orm import load_only, selectinload

from app.database import get_session
from app.models import ApiKey, Secret
from app.responses import text_response


def _qp(request: Request, name: str) -> str | None:
    raw = request.query_params.get(name, None)
    return unquote(raw) if raw is not None else None


def register_api_routes(app: Robyn):
    @app.get("/api/secret")
    async def get_secret_via_api(request: Request) -> Response:
        token = _qp(request, "token")
        key = _qp(request, "key")

        if not token:
            return text_response(401, "Unauthorized: missing or invalid token")
        if not key:
            return text_response(404, "Secret not found")

        with get_session() as session:
            secret = session.scalars(select(Secret).where(Secret.key == key)).first()
            if secret is None:
                return text_response(404, "Secret not found")

            api_key = session.scalars(select(ApiKey).options(selectinload(ApiKey.secrets).load_only(Secret.id)).where(ApiKey.key == token)).first()
            if api_key is None:
                return text_response(401, "Unauthorized: missing or invalid token")

            scoped_ids = {s.id for s in api_key.secrets}
            if scoped_ids and secret.id not in scoped_ids:
                return text_response(401, "Unauthorized: missing or invalid token")

            return text_response(200, secret.value)
