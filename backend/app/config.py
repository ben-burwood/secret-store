import os

PORT = int(os.environ.get("PORT", "8080"))
DB_PATH = os.environ.get("DB_PATH", "secret-store.db")

SESSION_COOKIE = "secret_store_session"

ENCRYPTION_KEY = os.environ.get("ENCRYPTION_KEY")
if not ENCRYPTION_KEY:
    raise RuntimeError("ENCRYPTION_KEY environment variable is required")

DASHBOARD_USER = os.environ.get("DASHBOARD_USER")
DASHBOARD_PASSWORD = os.environ.get("DASHBOARD_PASSWORD")
if not DASHBOARD_USER or not DASHBOARD_PASSWORD:
    raise RuntimeError("DASHBOARD_USER and DASHBOARD_PASSWORD environment variables are required")
