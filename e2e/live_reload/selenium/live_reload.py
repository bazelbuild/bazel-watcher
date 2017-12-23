import os
import sys
import socket

from SimpleHTTPServer import SimpleHTTPRequestHandler
from BaseHTTPServer import HTTPServer

path = "e2e/live_reload/test.txt"

class HttpHandler(SimpleHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header("Content-type", "text/html")
        self.end_headers()

        # Send some body content.
        self.wfile.write("""\
<html>
    <head>
        <script src="%s"></script>
    </head>
    <body>
        <p>%s</p>
    </body>
</html>
""" % (os.getenv("IBAZEL_LIVERELOAD_URL", ""), open(path).read()))

    def log_message(self, format, *args):
        pass

def main():
    """Brings up a http server on the supplied port which serves the local fs.
    """
    server = HTTPServer(('127.0.0.1', 0), HttpHandler)

    sys.stdout.write("Webserver url: http://127.0.0.1:%s/" % (server.server_port))
    sys.stdout.flush()

    server.serve_forever()

if __name__ == "__main__":
    try:
        open(path)
    except:
        path  = "test.txt"
    main()
