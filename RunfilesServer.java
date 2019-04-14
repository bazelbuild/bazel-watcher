package brs;

import static com.google.common.net.HttpHeaders.CONTENT_TYPE;
import static java.net.HttpURLConnection.HTTP_NOT_FOUND;
import static java.net.HttpURLConnection.HTTP_OK;

import com.beust.jcommander.JCommander;
import com.beust.jcommander.Parameter;
import com.google.common.base.Preconditions;
import com.google.common.flogger.FluentLogger;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;
import java.io.File;
import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.file.Files;
import java.util.concurrent.Executors;
import javax.activation.MimetypesFileTypeMap;

/** Simple web server that serves files directly out of the bazel runfiles tree. */
public final class RunfilesServer {

  private static final FluentLogger logger = FluentLogger.forEnclosingClass();
  private static final MimetypesFileTypeMap FILE_TYPE_MAP;

  static {
    try {
    // The mime.types included with the javax.activation jar is ancient. Use a mime.types from a
    // recent release of apache instead.
    FILE_TYPE_MAP =
        new MimetypesFileTypeMap("external/apache_mime_types/file/downloaded");
    } catch (IOException e) {
      throw new RuntimeException(e);
    }
  }

  @Parameter(names = "--port", required = true)
  private int port;

  private RunfilesServer() {}

  public static void main(String[] args) throws IOException {
    RunfilesServer me = new RunfilesServer();
    JCommander.newBuilder().addObject(me).build().parse(args);
    HttpServer server = HttpServer.create(new InetSocketAddress(me.port), 0 /* backlog */);
    server.setExecutor(Executors.newCachedThreadPool());
    server.createContext("/", RunfilesServer::handle);
    server.start();
    logger.atInfo().log("listening on port %d", me.port);
  }

  private static void handle(HttpExchange httpExchange) throws IOException {
    String path = httpExchange.getRequestURI().toString();
    Preconditions.checkState(path.startsWith("/"));
    String runfilesPath = path.substring(1);
    File runfile = new File(runfilesPath);
    int status;
    if (!runfile.exists()) {
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
}
