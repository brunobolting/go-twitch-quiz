import config from "./config.js"
import Game from "./game.js"

const modal = document.querySelector(".modal-error")
const error = modal.querySelector(".error")

const run = () => {
    const searchParams = new URLSearchParams(window.location.search);
    const channel = searchParams.get("chat-from");
    console.log(channel)
    if (channel === "" || channel === null) {
        modal.style.display = "block";
        error.innerHTML = "Chat não informado..."
        return
    }

    const conn = new WebSocket(`${config.WS_URL}?channel=${channel}`)
    const game = new Game(conn)

    conn.onmessage = (event) => {
        game.handleMessage(event)
    }

    conn.onclose = (event) => {
        setTimeout((event) => {

            game.stop()
            modal.style.display = "block";
            error.innerHTML = "A conexão com o servidor foi perdida..."
        }, 2000)
    };

    conn.addEventListener("open", (event) => {
        game.startRound()
    })
}

window.addEventListener("load", () => {
    run()
})
