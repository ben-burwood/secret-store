import json
from datetime import datetime

from robyn import Request, Response, Robyn
from sqlalchemy import delete, select

from app.database import get_session
from app.models import Secret
from app.responses import empty, form_field, json_response, parse_json_body


def _serialize(secret: Secret) -> dict:
    return {
        "id": secret.id,
        "key": secret.key,
        "value": secret.value,
        "created_at": secret.created_at.isoformat() if secret.created_at else None,
    }


def _parse_iso(value) -> datetime | None:
    if not isinstance(value, str):
        return None
    try:
        return datetime.fromisoformat(value.replace("Z", "+00:00"))
    except ValueError:
        return None


def _build_secrets(items: list, *, default_created_at: datetime, preserve_created_at: bool) -> list[Secret]:
    out: list[Secret] = []
    for item in items:
        key = (item.get("key") or "").strip()
        if not key:
            continue
        created_at = (preserve_created_at and _parse_iso(item.get("created_at"))) or default_created_at
        out.append(Secret(key=key, value=item.get("value") or "", created_at=created_at))
    return out


def _read_secrets_form(request: Request) -> tuple[list | None, Response | None]:
    raw = form_field(request, "secrets")
    if raw is None:
        return None, json_response(400, {"error": "Missing 'secrets' field"})
    try:
        payload = json.loads(raw)
    except json.JSONDecodeError:
        return None, json_response(400, {"error": "Invalid JSON in 'secrets'"})
    items = payload if isinstance(payload, list) else payload.get("secrets", [])
    if not isinstance(items, list):
        return None, json_response(400, {"error": "'secrets' must be a list"})
    return items, None


def register_secrets_routes(app: Robyn):
    @app.get("/web/secrets", auth_required=True)
    async def list_secrets(request: Request) -> Response:
        with get_session() as session:
            rows = session.scalars(select(Secret).order_by(Secret.id)).all()
            return json_response(200, {"secrets": [_serialize(s) for s in rows]})

    @app.post("/web/secrets/new", auth_required=True)
    async def create_secret(request: Request) -> Response:
        body, err = parse_json_body(request)
        if err is not None:
            return err

        key = (body.get("key") or "").strip()
        if not key:
            return json_response(400, {"error": "key is required"})

        with get_session() as session:
            session.add(Secret(key=key, value=body.get("value") or "", created_at=datetime.utcnow()))
            session.commit()
        return empty(201)

    @app.patch("/web/secrets/:id", auth_required=True)
    async def update_secret(request: Request) -> Response:
        try:
            secret_id = int(request.path_params["id"])
        except (KeyError, ValueError):
            return json_response(400, {"error": "Invalid ID"})

        body, err = parse_json_body(request)
        if err is not None:
            return err

        with get_session() as session:
            secret = session.get(Secret, secret_id)
            if secret is None:
                return json_response(404, {"error": "Secret not found"})
            if "key" in body:
                secret.key = body["key"]
            if "value" in body:
                secret.value = body["value"]
            secret.created_at = datetime.utcnow()
            session.commit()
        return empty(204)

    @app.delete("/web/secrets/:id", auth_required=True)
    async def delete_secret(request: Request) -> Response:
        try:
            secret_id = int(request.path_params["id"])
        except (KeyError, ValueError):
            return json_response(400, {"error": "Invalid ID"})

        with get_session() as session:
            session.execute(delete(Secret).where(Secret.id == secret_id))
            session.commit()
        return empty(204)

    @app.post("/web/import", auth_required=True)
    async def import_secrets(request: Request) -> Response:
        items, err = _read_secrets_form(request)
        if err is not None:
            return err

        with get_session() as session:
            session.add_all(_build_secrets(items, default_created_at=datetime.utcnow(), preserve_created_at=False))
            session.commit()
        return empty(204)

    @app.post("/web/restore", auth_required=True)
    async def restore_secrets(request: Request) -> Response:
        items, err = _read_secrets_form(request)
        if err is not None:
            return err

        with get_session() as session:
            session.execute(delete(Secret))
            session.add_all(_build_secrets(items, default_created_at=datetime.utcnow(), preserve_created_at=True))
            session.commit()
        return empty(204)
