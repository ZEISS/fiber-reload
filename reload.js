(() => {
  // reload.ts
  var lastUuid = "";
  var timeout;
  var resetBackoff = () => {
    timeout = 1e3;
  };
  var backOff = () => {
    if (timeout > 10 * 1e3) {
      return;
    }
    timeout = timeout * 2;
  };
  var hotReloadUrl = () => {
    const hostAndPort = location.hostname + (location.port ? ":" + location.port : "");
    if (location.protocol === "https:") {
      return "wss://" + hostAndPort + "/ws/reload";
    }
    return "ws://" + hostAndPort + "/ws/reload";
  };
  function connectHotReload() {
    const socket = new WebSocket(hotReloadUrl());
    socket.onmessage = (event) => {
      if (lastUuid === "") {
        lastUuid = event.data;
      }
      if (lastUuid !== event.data) {
        console.log("[Hot Reloader] Server Changed, reloading");
        location.reload();
      }
    };
    socket.onopen = () => {
      resetBackoff();
      socket.send("Hello");
    };
    socket.onclose = () => {
      const timeoutId = setTimeout(function() {
        clearTimeout(timeoutId);
        backOff();
        connectHotReload();
      }, timeout);
    };
  }
  resetBackoff();
  connectHotReload();
})();
