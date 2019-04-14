package brs;

import java.io.IOException;
import java.net.ServerSocket;

public final class EphemeralPort {

  private EphemeralPort() {}

  public static int get() throws IOException {
    // The ServerSocket ctor both chooses an ephemeral port and automatically binds to it.
    // We need to close the socket to prevent a BindException from being thrown from
    // HttpServer.create. This creates a small race condition -- something else could bind to the
    // ephemeral port in the meantime.
    // See https://stackoverflow.com/questions/2675362/how-to-find-an-available-port.
    try (ServerSocket whatever = new ServerSocket(0)) {
      return whatever.getLocalPort();
    }
  }
}
