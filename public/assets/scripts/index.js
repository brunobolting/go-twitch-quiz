// import Ws from "./ws.js"
import Game from "./game.js"

const WS_URL = "ws://localhost:8080/ws"

const run = () => {
    const searchParams = new URLSearchParams(window.location.search);
    const channel = searchParams.get("chat-from");

    const conn = new WebSocket(`${WS_URL}?channel=${channel}`)
    const game = new Game(conn)

    conn.onmessage = (event) => {
        game.handleMessage(event)
    }

    conn.onclose = (event) => {
        console.log("connection closed", event)
    };

    conn.addEventListener("open", (event) => {
        game.startRound()
    })

    // setTimeout(() => {
    //     console.log("one second later...")
    //     conn.send(JSON.stringify({command: "START_ROUND"}))
    // }, 1000)
}

window.addEventListener("load", () => {
    console.log('aoppppa')
    run()
})
