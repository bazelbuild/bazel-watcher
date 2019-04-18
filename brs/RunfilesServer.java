package brs;

import static com.google.common.net.HttpHeaders.CONTENT_TYPE;
import static java.net.HttpURLConnection.HTTP_NOT_FOUND;
import static java.net.HttpURLConnection.HTTP_OK;
import static java.nio.charset.StandardCharsets.UTF_8;

import com.beust.jcommander.JCommander;
import com.beust.jcommander.Parameter;
import com.google.common.flogger.FluentLogger;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;
import java.awt.Desktop;
import java.io.File;
import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.net.URI;
import java.net.URISyntaxException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.concurrent.Executors;
import javax.activation.MimetypesFileTypeMap;
import javax.annotation.Nullable;

/**
 * Simple web server that serves files directly out of a bazel target's runfiles.
 * Powers the {@code serve} rule and the {@code serve_this} Starlark helper function.
 *
 * The goal of this server is to serve runfiles with no modification (bundling, minification,
 * etc.) other than what is required for ibazel integration. Specifically, when this server is run
 * under ibazel, it will inject script tags pointing to the ibazel livereload snippet into every
 * HTML payload.
 */
public final class RunfilesServer {

  @Parameter(
      names = "--port",
      description = "port to listen on. If not given, an ephemeral port will be chosen")
  private int port;

  @Parameter(
      names = "--index",
      description =
          "page to visit in the system's default browser when the server is up. If not given, the "
              + "browser will not be launched.")
  private String indexToOpen;

  @Parameter(
      names = "--nobrowser",
      description =
          "Disables opening the browser."
              + "The default behavior of RunfilesServer is to open a browser to the page given by --index. "
              + "Pass --nobrowser if this behavior is not appropriate. "
              + "Example: bazel run //some_serve_target -- --nobrowser")
  private boolean disableBrowser;

  private static final FluentLogger logger = FluentLogger.forEnclosingClass();
  private static final MimetypesFileTypeMap FILE_TYPE_MAP;
  private static final Path CWD = Paths.get("").toAbsolutePath();
  @Nullable private static final byte[] LIVERELOAD_SNIPPET;

  static {
    // If the this server is being run under ibazel as part of a target that has the tag
    // `ibazel_live_reload`, ibazel will set the IBAZEL_LIVERELOAD_URL environment variable.
    String livereloadUrl = System.getenv("IBAZEL_LIVERELOAD_URL");
    LIVERELOAD_SNIPPET =
        livereloadUrl == null
            ? null
            : String.format("<script src=\"%s\"></script>", livereloadUrl).getBytes(UTF_8);
  }

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

  public static void main(String[] args) throws IOException, URISyntaxException {
    RunfilesServer me = new RunfilesServer();
    JCommander.newBuilder().addObject(me).build().parse(args);
    int port = me.port == 0 ? EphemeralPort.get() : me.port;
    HttpServer server = HttpServer.create(new InetSocketAddress(port), 0 /* backlog */);
    server.setExecutor(Executors.newCachedThreadPool());
    server.createContext("/", RunfilesServer::handle);
    server.start();
    logger.atInfo().log("listening on port %d", port);
    // Print a line to stdout. IntegrationTestRunner uses this for synchronization (it won't run the
    // test binary until the system under test prints a line to stdout). For other uses, this is
    // harmless.
    System.out.println("ok");
    if (me.shouldOpenBrowser()) {
      String index = String.format("http://localhost:%d/%s", port, me.indexToOpen);
      try {
        // The getDesktop api apparently has never worked on linux:
        // https://stackoverflow.com/questions/18004150/desktop-api-is-not-supported-on-the-current-platform
        // Fall back to an X Windows api for best effort.
        // Write once, debug everywhere...
        Desktop.getDesktop().browse(new URI(index));
      } catch (UnsupportedOperationException e) {
        new ProcessBuilder("xdg-open", index).start();
      }
    }
  }

  private boolean shouldOpenBrowser() {
    return indexToOpen != null && !disableBrowser;
  }

  private static void handle(HttpExchange httpExchange) throws IOException {
    String path = httpExchange.getRequestURI().toString();
    @Nullable File runfile = resolve(path);
    int status;
    if (runfile == null || !runfile.exists()) {
      httpExchange.sendResponseHeaders(status = HTTP_NOT_FOUND, 0);
      httpExchange.getResponseBody().close();
    } else {
      status = sendFile(httpExchange, runfile);
    }
    logger.atInfo().log("%d %s", status, path);
  }

  private static int sendFile(HttpExchange httpExchange, File file) throws IOException {
    int status = HTTP_OK;
    String contentType = FILE_TYPE_MAP.getContentType(file);
    // TODO: consider not hardcoding this.
    if (contentType.startsWith("text/")) {
      contentType = contentType + "; charset=utf-8";
    }
    httpExchange.getResponseHeaders().add(CONTENT_TYPE, contentType);

    long length = file.length();
    boolean shouldInjectLiveReloadSnippet =
        LIVERELOAD_SNIPPET != null && file.getPath().endsWith(".html");
    if (shouldInjectLiveReloadSnippet) {
      length += LIVERELOAD_SNIPPET.length;
    }
    httpExchange.sendResponseHeaders(status, length);
    try (OutputStream out = httpExchange.getResponseBody()) {
      Files.copy(file.toPath(), out);
      if (shouldInjectLiveReloadSnippet) {
        out.write(LIVERELOAD_SNIPPET);
      }
    }
    return status;
  }

  @Nullable
  private static File resolve(String path) {
    // Paths derived by a user agent from a url should be absolute. Guard against handcrafted paths.
    if (!path.startsWith("/")) {
      return null;
    }
    // Don't allow traversing out of the runfiles. Based on https://stackoverflow.com/a/33084369.
    // IMPORTANT: RunfilesServer is designed for local development (and it is not clear how it would
    // be used in production, since the concept of runfiles and bazel itself should not exist
    // there). This check helps developers avoid accidentally serving files in a nonhermetic way;
    // it may be not enough to guard against malicious path-traversal attacks
    // (see https://en.wikipedia.org/wiki/Directory_traversal_attack). In particular, it does not
    // use a chroot.
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
