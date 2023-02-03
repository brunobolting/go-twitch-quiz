class Ws {
    constructor(url) {
        this.ws = new WebSocket(url)
    }

    Send(data) {
        this.ws.send(JSON.stringify(data))
    }

    OnMessage(callable) {
        this.ws.onmessage = callable
    }

    OnClose(callable) {
        this.ws.onclose = callable
    }

    GetConnection() {
        return this.ws
    }
}

export default Ws
