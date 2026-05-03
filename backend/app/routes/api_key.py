import base64
import secrets as pysecrets
from datetime import datetime

from robyn import Request, Response, Robyn
from sqlalchemy import delete, select
from sqlalchemy.exc import IntegrityError
from sqlalchemy.orm import selectinload

from app.database import get_session
from app.models import ApiKey, Secret
from app.responses import empty, json_response, parse_json_body


def _new_key() -> str:
    return base64.b64encode(pysecrets.token_bytes(32)).decode("ascii")


def _serialize(api_key: ApiKey) -> dict:
    return {
        "id": api_key.id,
        "name": api_key.name,
        "key": api_key.key,
        "created_at": api_key.created_at.isoformat() if api_key.created_at else None,
        "secret_ids": sorted(s.id for s in api_key.secrets),
    }


def register_api_key_routes(app: Robyn):
    @app.get("/web/api/keys", auth_required=True)
    async def list_api_keys(request: Request) -> Response:
        with get_session() as session:
            rows = session.scalars(select(ApiKey).options(selectinload(ApiKey.secrets)).order_by(ApiKey.id)).all()
            return json_response(200, {"api_keys": [_serialize(k) for k in rows]})

    @app.post("/web/api/keys/new", auth_required=True)
    async def create_api_key(request: Request) -> Response:
        body, err = parse_json_body(request)
        if err is not None:
            return err

        name = (body.get("name") or "").strip()
        if not name:
            return json_response(400, {"error": "name is required"})

        with get_session() as session:
            row = ApiKey(name=name, key=_new_key(), created_at=datetime.now())
            session.add(row)
            try:
                session.commit()
            except IntegrityError:
                return json_response(409, {"error": f"Name '{name}' already exists"})
            return json_response(201, _serialize(row))

    @app.post("/web/api/keys/:id/regenerate", auth_required=True)
    async def regenerate_api_key(request: Request) -> Response:
        try:
            api_key_id = int(request.path_params["id"])
        except (KeyError, ValueError):
            return json_response(400, {"error": "Invalid ID"})

        with get_session() as session:
            row = session.get(ApiKey, api_key_id)
            if row is None:
                return json_response(404, {"error": "API key not found"})
            row.key = _new_key()
            session.commit()
            return json_response(200, _serialize(row))

    @app.put("/web/api/keys/:id/scopes", auth_required=True)
    async def update_api_key_scopes(request: Request) -> Response:
        try:
            api_key_id = int(request.path_params["id"])
        except (KeyError, ValueError):
            return json_response(400, {"error": "Invalid ID"})

        body, err = parse_json_body(request)
        if err is not None:
            return err

        raw_ids = body.get("secret_ids") if isinstance(body, dict) else None
        if not isinstance(raw_ids, list) or not all(isinstance(x, int) for x in raw_ids):
            return json_response(400, {"error": "secret_ids must be a list of integers"})
        secret_ids = list(set(raw_ids))

        with get_session() as session:
            row = session.get(ApiKey, api_key_id)
            if row is None:
                return json_response(404, {"error": "API key not found"})
            if secret_ids:
                matched = session.scalars(select(Secret).where(Secret.id.in_(secret_ids))).all()
                if len(matched) != len(secret_ids):
                    return json_response(400, {"error": "One or more secret_ids do not exist"})
                row.secrets = list(matched)
            else:
                row.secrets = []
            session.commit()
            return json_response(200, _serialize(row))

    @app.delete("/web/api/keys/:id", auth_required=True)
    async def delete_api_key(request: Request) -> Response:
        try:
            api_key_id = int(request.path_params["id"])
        except (KeyError, ValueError):
            return json_response(400, {"error": "Invalid ID"})

        with get_session() as session:
            session.execute(delete(ApiKey).where(ApiKey.id == api_key_id))
            session.commit()
        return empty(204)
