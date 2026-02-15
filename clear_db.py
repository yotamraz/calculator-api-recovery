"""Clear all calculation records from the database."""
import sqlite3
import pathlib

db_path = pathlib.Path("calculator.db")
if db_path.exists():
    conn = sqlite3.connect(str(db_path))
    try:
        conn.execute("DELETE FROM calculation")
        conn.commit()
    except Exception:
        pass
    finally:
        conn.close()
