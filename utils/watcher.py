import subprocess
import time
import threading
from datetime import datetime, timedelta
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler

class Watcher(FileSystemEventHandler):
    def __init__(self, path, backend=None, on_modified=None, on_created=None, on_deleted=None, on_moved=None, on_any=None):
        self.path = path
        self.backend = backend
        self.last_modified = datetime.now()

        if backend == "watchdog":
            self.observer = Observer()
            self.observer.schedule(self, self.path, recursive=False)

            if on_any == None:
                self.handle_on_created = on_created
                self.handle_on_modified = on_modified
                self.handle_on_deleted = on_deleted
                self.handle_on_moved = on_moved
            else:
                self.handle_on_created = on_any
                self.handle_on_modified = on_any
                self.handle_on_deleted = on_any
                self.handle_on_moved = on_any

        return


    def start(self):
        match self.backend:
            case "watchdog":
                def do():
                    self.observer.start()
                    while True:
                        time.sleep(1)

                background = threading.Thread(target=do, args=())
                background.daemon = True
                background.start()
            case "air":
                cmd = subprocess.Popen("air", cwd=self.path, shell=True)
                if cmd.wait() != 0:
                    raise ValueError(f"air returned non 0 status: {cod}")
            case _:
                raise ValueError("Invalid backend value")

        return
    

    def maybe_block(self):
        if datetime.now() - self.last_modified < timedelta(seconds=1):
            return

        self.last_modified = datetime.now()

        return


    def on_created(self, event):
        self.maybe_block()
        self.handle_on_created(self, event)

        return


    def on_modified(self, event):
        self.maybe_block()
        self.handle_on_modified(self, event)

        return


    def on_deleted(self, event):
        self.maybe_block()
        self.handle_on_deleted(self, event)

        return


    def on_moved(self, event):
        self.maybe_block()
        self.handle_on_moved(self, event)

        return
