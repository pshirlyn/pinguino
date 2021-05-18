"use strict";

document.addEventListener('DOMContentLoaded', DOMloaded, false);

function Session() {
    this.socket = null;
    this.socketOpen = false;
    this.terminated = false;
}

Session.prototype.sendMessage = function (jsonobject) {
    if (this.socketOpen) {
        this.socket.send(JSON.stringify(jsonobject));
    }
};

Session.prototype.onclose = function () {
    console.log("socket closed!");
};

Session.prototype.reportClose = function () {
    if (!this.terminated) {
        this.terminated = true;
        if (this.onclose !== null) {
            this.onclose();
        }
    }
};

Session.prototype.openSocket = function () {
    if (this.socket !== null) {
        throw new Error("socket already created for session");
    }
    const url = new URL("/ws", window.location.href);
    url.protocol = (url.protocol === "http:") ? "ws:" : "wss:";
    // let socket = new WebSocket("ws://127.0.0.1:8080/ws");
    console.log("connecting to", url.href);
    this.socket = new WebSocket(url.href);
};

Session.prototype.onMessage = function(msg) {
    console.log("received: ", msg);
    // this.sendMessage("client replies hello!");
};

Session.prototype.connect = function () {
    const session = this;
    this.openSocket();
    this.socketOpen = false;
    this.terminated = false;
    this.socket.addEventListener("open", function () {
        console.log("connection opened");
        session.socketOpen = true;
    });
    this.socket.addEventListener('error', function () {
        console.log("connection error");
        session.reportClose();
    });
    this.socket.addEventListener('message', function (ev) {
        console.log("client received message", ev.data);
        session.onMessage(ev.data);
        
    });
    this.socket.addEventListener('close', function () {
        console.log("connection terminated");
        session.reportClose();
    });
};

function Canvas(canvas) {
    this.canvas = canvas;
    this.viewWidth = this.viewHeight = 100;
    this.width = this.height = 100;
    this.gameSprites = [];
}

Canvas.prototype.testDraw = function() {
    this.canvas.moveTo(0, 0);
    this.canvas.lineTo(200, 100);
    this.canvas.stroke();
}

Canvas.prototype.updateSprites = function (sprites) {
    this.gameSprites = sprites;
};

Canvas.prototype.getMousePosition = function (ev) {
    const rect = this.canvas.getBoundingClientRect();
    let x = ev.clientX - rect.left;
    let y = ev.clientY - rect.top;
    return {x: x, y: y};
};


function DOMloaded() {

    var session = new Session();
    session.connect();
    

    var c = document.getElementById("main-canvas");
    var ctx = c.getContext("2d");
    
    var canvas = new Canvas(ctx);
    canvas.testDraw()

    var msg = document.getElementById("msg");
    var log = document.getElementById("log");

    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }
    

    document.getElementById("form").onsubmit = function () {
        if (!session) {
            return false;
        }
        if (!msg.value) {
            return false;
        }
        
        session.sendMessage(msg.value);
        msg.value = "";
        return false;
    };
}

