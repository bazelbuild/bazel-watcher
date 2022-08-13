// Copyright 2017 The Bazel Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

(function(){
  function extractUrl() {
    return [...window.document.getElementsByTagName('script') ]
        .filter(function(e) { return e.src; })
        .map(function(e) {
          const a = document.createElement('a');
          a.href = e.src;
          return a;
        })
        .filter(function (a) { return a.pathname == "/profiler.js"; })
        .map(function (a) { return `${a.protocol}//${a.host}/profiler-event`; })
        .reduce(function(acc, u) { return acc || u; });
  };

  const profilerEventUrl = extractUrl();
  if (!profilerEventUrl) {
    this.console.error("iBazel profiler disabled because it could not find its own <SCRIPT> tag");
    return;
  }

  if (!window.navigator.sendBeacon) {
    console.error(
        "iBazel profiler disabled because Beacon API is not available")
    return;
  }

  function profilerEvent(eventType, eventData) {
    if (!eventType || typeof eventType != 'string') {
      console.error("Invalid iBazel profiler event type", eventType);
      return;
    }
    if (typeof eventData == 'object') {
      eventData = JSON.stringify(eventData)
    }
    const event = { type: eventType, time: Date.now() };
    event.timeSinceNavigationStart = event.time - window.performance.timing.navigationStart;
    if (eventData) {
      event.data = eventData.toString();
    }
    const data = JSON.stringify(event)
    if (!window.navigator.sendBeacon(profilerEventUrl, data)) {
      console.error("Failed to send profile data to ", profilerEventUrl);
    }
  }

  // Expose a public window.IBazelProfileEvent API for the user space to be
  // able to send profile events
  window.IBazelProfileEvent = profilerEvent;

  window.addEventListener("load", function() {
    let timing = window.performance.timing;

    profilerEvent("PAGE_LOAD", {
      // deltas
      pageLoadTime: timing.loadEventStart - timing.navigationStart, // loadEventEnd is not set yet
      fetchTime: timing.responseEnd - timing.fetchStart,
      connectTime: timing.connectEnd - timing.connectStart,
      requestTime: timing.responseEnd - timing.requestStart,
      responseTime: timing.responseEnd - timing.responseStart,
      renderTime: timing.domComplete - timing.domLoading,

      // absolutes
      navigationStart: timing.navigationStart,
      unloadEventStart: timing.unloadEventStart,
      unloadEventEnd: timing.unloadEventEnd,
      redirectStart: timing.redirectStart,
      redirectEnd: timing.redirectEnd,
      fetchStart: timing.fetchStart,
      domainLookupStart: timing.domainLookupStart,
      domainLookupEnd: timing.domainLookupEnd,
      connectStart: timing.connectStart,
      connectEnd: timing.connectEnd,
      secureConnectionStart: timing.secureConnectionStart,
      requestStart: timing.requestStart,
      responseStart: timing.responseStart,
      responseEnd: timing.responseEnd,
      domLoading: timing.domLoading,
      domInteractive: timing.domInteractive,
      domContentLoadedEventStart: timing.domContentLoadedEventStart,
      domContentLoadedEventEnd: timing.domContentLoadedEventEnd,
      domComplete: timing.domComplete,
      loadEventStart: timing.loadEventStart,
    });
  }, false);
})();