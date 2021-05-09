"use strict";

/* Session handles an open websocket connection to the server.
 * The public interface:
 *  - sendMessage(jsonobject)
 *  - connect()
 * The events to be overridden by the user:
 *  - onmessage
 *  - onclose
 */

let socket = new WebSocket("ws://127.0.0.1:8080/ws");
console.log("Attempting Connection...");

socket.onopen = () => {
    console.log("Successfully Connected");
    socket.send("Hi From the Client!")
};

socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
    socket.send("Client Closed!")
};

socket.onerror = error => {
    console.log("Socket Error: ", error);
};

// function Session() {
//     this.socket = null;
//     this.onmessage = null;
//     this.onclose = null;
//     this.socketOpen = false;
//     this.terminated = false;
// }

// Session.prototype.sendMessage = function (jsonobject) {
//     if (this.socketOpen) {
//         this.socket.send(JSON.stringify(jsonobject));
//     }
// };

// Session.prototype.reportClose = function () {
//     if (!this.terminated) {
//         this.terminated = true;
//         if (this.onclose !== null) {
//             this.onclose();
//         }
//     }
// };

// Session.prototype.openSocket = function () {
//     if (this.socket !== null) {
//         throw new Error("socket already created for session");
//     }
//     const url = new URL("/websocket", window.location.href);
//     url.protocol = (url.protocol === "http:") ? "ws:" : "wss:";
//     console.log("connecting to", url.href);
//     this.socket = new WebSocket(url.href);
// };

// Session.prototype.connect = function () {
//     const session = this;
//     this.openSocket();
//     this.socketOpen = false;
//     this.terminated = false;
//     this.socket.addEventListener("open", function () {
//         console.log("connection opened");
//         session.socketOpen = true;
//     });
//     this.socket.addEventListener('error', function () {
//         console.log("connection error");
//         session.reportClose();
//     });
//     this.socket.addEventListener('message', function (ev) {
//         if (session.onmessage !== null) {
//             session.onmessage(JSON.parse(ev.data));
//         }
//     });
//     this.socket.addEventListener('close', function () {
//         console.log("connection terminated");
//         session.reportClose();
//     });
// };


// function Canvas(canvas) {
//     this.canvas = canvas;
//     this.viewWidth = this.viewHeight = 100;
//     this.width = this.height = 100;
//     this.gameSprites = [];
// }


// Canvas.prototype.updateSprites = function (sprites) {
//     this.gameSprites = sprites;
// };

// Canvas.prototype.getMousePosition = function (ev) {
//     const rect = this.canvas.getBoundingClientRect();
//     let x = ev.clientX - rect.left;
//     let y = ev.clientY - rect.top;
//     return {x: x, y: y};
// };
