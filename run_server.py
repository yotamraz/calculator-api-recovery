"""Helper script to start the server."""
import uvicorn
from server import app, settings

uvicorn.run(app, host=settings.server_host, port=settings.server_port)
