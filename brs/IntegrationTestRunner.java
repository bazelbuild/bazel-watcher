package brs;

import static com.google.common.base.Preconditions.checkState;
import static java.lang.ProcessBuilder.Redirect.INHERIT;

import com.beust.jcommander.JCommander;
import com.beust.jcommander.Parameter;
import com.google.common.flogger.FluentLogger;
import com.google.common.util.concurrent.FutureCallback;
import com.google.common.util.concurrent.Futures;
import com.google.common.util.concurrent.MoreExecutors;
import com.google.common.util.concurrent.SettableFuture;
import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.util.concurrent.CancellationException;

public final class IntegrationTestRunner {

  private static final FluentLogger logger = FluentLogger.forEnclosingClass();

  @Parameter(names = "--sut_binary", required = true)
  private String sutBinary;

  @Parameter(names = "--test_binary", required = true)
  private String testBinary;

  private IntegrationTestRunner() {}

  public static void main(String[] args) throws IOException {
    IntegrationTestRunner runner = new IntegrationTestRunner();
    JCommander.newBuilder().addObject(runner).build().parse(args);
    runner.run();
  }

  private void run() throws IOException {
    String port = Integer.toString(EphemeralPort.get());

    // Bring up the system under test.
    Process systemUnderTest =
        new ProcessBuilder(sutBinary, "--port", port).redirectError(INHERIT).start();
    checkState(systemUnderTest.isAlive(), "%s already died!", sutBinary);
    // Block until the system under test writes a line to its stdout.
    new BufferedReader(new InputStreamReader(systemUnderTest.getInputStream())).readLine();

    // Run the test binary.
    Process testProcess =
        new ProcessBuilder(testBinary, "--backend_port", port)
            .redirectError(INHERIT)
            .redirectOutput(INHERIT)
            .start();

    SettableFuture<Integer> testStatus = SettableFuture.create();
    new Thread(
            () -> {
              int status;
              try {
                status = testProcess.waitFor();
              } catch (InterruptedException e) {
                testStatus.setException(e);
                return;
              }
              testStatus.set(status);
            })
        .start();

    Futures.addCallback(
        testStatus,
        new FutureCallback<>() {
          @Override
          public void onSuccess(Integer result) {
            logger.atInfo().log("test binary %s exited with status %d", testBinary, result);
            System.exit(result);
          }

          @Override
          public void onFailure(Throwable throwable) {
            checkState(throwable instanceof CancellationException);
            System.exit(1);
          }
        },
        MoreExecutors.directExecutor());
  }
}
