import json

from robyn import Request, Response

_JSON = {"content-type": "application/json"}
_TEXT = {"content-type": "text/plain"}


def json_response(status_code: int, payload) -> Response:
    return Response(status_code=status_code, headers=_JSON, description=json.dumps(payload))


def empty(status_code: int) -> Response:
    return Response(status_code=status_code, headers=_JSON, description="")


def text_response(status_code: int, body: str) -> Response:
    return Response(status_code=status_code, headers=_TEXT, description=body)


def parse_json_body(request: Request) -> tuple[object, Response | None]:
    try:
        return json.loads(request.body), None
    except (json.JSONDecodeError, TypeError):
        return None, json_response(400, {"error": "Invalid JSON body"})


def form_field(request: Request, name: str) -> str | None:
    raw = request.form_data.get(name) if request.form_data else None
    if raw is None:
        return None
    if isinstance(raw, bytes):
        return raw.decode("utf-8")
    return raw
