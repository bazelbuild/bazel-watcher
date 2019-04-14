package brs;

import static com.google.common.net.HttpHeaders.CONTENT_TYPE;
import static java.net.HttpURLConnection.HTTP_NOT_FOUND;
import static java.net.HttpURLConnection.HTTP_OK;

import com.beust.jcommander.JCommander;
import com.beust.jcommander.Parameter;
import com.google.common.flogger.FluentLogger;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;
import java.io.File;
import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.concurrent.Executors;
import javax.activation.MimetypesFileTypeMap;
import javax.annotation.Nullable;

/** Simple web server that serves files directly out of the bazel runfiles tree. */
public final class RunfilesServer {

  @Parameter(
      names = "--port",
      description = "port to listen on. If not given, an ephemeral port will be chosen")
  private int port;

  private static final FluentLogger logger = FluentLogger.forEnclosingClass();
  private static final MimetypesFileTypeMap FILE_TYPE_MAP;
  private static final Path CWD = Paths.get("").toAbsolutePath();

  static {
    try {
      // The mime.types included with the javax.activation jar is ancient. Use a mime.types from a
      // recent release of apache instead.
      FILE_TYPE_MAP = new MimetypesFileTypeMap("external/apache_mime_types/file/downloaded");
    } catch (IOException e) {
      throw new RuntimeException(e);
    }
  }

  private RunfilesServer() {}

  public static void main(String[] args) throws IOException {
    RunfilesServer me = new RunfilesServer();
    JCommander.newBuilder().addObject(me).build().parse(args);
    int port = me.port == 0 ? EphemeralPort.get() : me.port;
    HttpServer server = HttpServer.create(new InetSocketAddress(port), 0 /* backlog */);
    server.setExecutor(Executors.newCachedThreadPool());
    server.createContext("/", RunfilesServer::handle);
    server.start();
    logger.atInfo().log("listening on port %d", port);
    System.out.println("ok");
  }

  private static void handle(HttpExchange httpExchange) throws IOException {
    String path = httpExchange.getRequestURI().toString();
    @Nullable File runfile = resolve(path);
    int status;
    if (runfile == null || !runfile.exists()) {
      httpExchange.sendResponseHeaders(status = HTTP_NOT_FOUND, 0);
      httpExchange.getResponseBody().close();
    } else {
      httpExchange.getResponseHeaders().add(CONTENT_TYPE, FILE_TYPE_MAP.getContentType(runfile));
      httpExchange.sendResponseHeaders(status = HTTP_OK, runfile.length());
      try (OutputStream out = httpExchange.getResponseBody()) {
        Files.copy(runfile.toPath(), out);
      }
    }
    logger.atInfo().log("%d %s", status, path);
  }

  @Nullable
  private static File resolve(String path) {
    // Paths derived by a user agent from a url should be absolute. Guard against handcrafted paths.
    if (!path.startsWith("/")) {
      return null;
    }
    // Don't allow path traversal. Based on https://stackoverflow.com/a/33084369.
    Path untrusted = Paths.get(path.substring(1));
    if (untrusted.isAbsolute()) {
      return null;
    }
    Path resolved = CWD.resolve(untrusted).normalize();
    if (!resolved.startsWith(CWD)) {
      return null;
    }
    return resolved.toFile();
  }
}
