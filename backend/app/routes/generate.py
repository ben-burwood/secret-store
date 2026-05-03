import secrets as pysecrets
import string

from robyn import Request, SubRouter

from app.responses import json_response

_SYMBOLS = "!@#$%^&*()-_=+[]{}|;:,.<>?/"

router = SubRouter(__name__, prefix="/web")


def _truthy(value: str | None) -> bool:
    return (value or "").lower() == "true"


@router.get("/secret/generate")
async def generate_secret(request: Request):
    length_raw = request.query_params.get("length", None)
    try:
        length = int(length_raw) if length_raw is not None else 0
    except ValueError:
        return json_response(400, {"error": "Invalid length"})
    if length < 8:
        return json_response(400, {"error": "Invalid length"})

    charset = string.ascii_letters
    if _truthy(request.query_params.get("includeNumbers", None)):
        charset += string.digits
    if _truthy(request.query_params.get("includeSymbols", None)):
        charset += _SYMBOLS

    secret = "".join(pysecrets.choice(charset) for _ in range(length))

    return {"secret": secret}
