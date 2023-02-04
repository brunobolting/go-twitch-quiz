import CountDown from "./countdown.js"

class Game {
    constructor(conn) {
        this.conn = conn
        this.timer = document.querySelector(".time-box > .timer")
        this.question = document.querySelector(".question-box")
        this.timerTitle = document.querySelector(".alert > .title")
        this.timerMessage = document.querySelector(".alert > .message")
        this.topPage = document.querySelector(".container > .top-page")
        this.winnerWrap = document.querySelector(".answer-box")
        this.winner = document.querySelector(".answer-box > .user")
        this.answer = document.querySelector(".answer-box > .answer")
        this.waitBar = document.querySelector(".answer-wrap > .waiting")
        this.leaderboard = document.querySelector(".winners")
        this.roundTime = 81
        this.breaktimeTime = 16
        this.winnersList = []
    }

    startRound() {
        this.conn.send(JSON.stringify({command: "START_ROUND"}))

        this.topPage.classList.remove("win-background", "loose-background")
        this.topPage.classList.add("default-background")

        this.waitBar.classList.remove("hidden")
        this.winnerWrap.classList.add("hidden")
        this.timerTitle.innerHTML = "Acerte a resposta!"
        this.timerMessage.innerHTML = "Digite seus palpites no chat!"

        const render = (time) => {
            let seconds = time.totalInSeconds
            if (seconds < 10) {
                seconds = this.formatSeconds(seconds)
            }
            if (seconds <= 5) {
                seconds = `<spam class="color-red">${seconds}</spam>`
            }
            this.timer.innerHTML = seconds
        }
        const complete = () => {
            this.looseScreen()
        }

        this.countdown = new CountDown(this.getRoundTimer(this.roundTime), render, complete)
    }

    handleMessage(event) {
        const message = JSON.parse(event.data)

        if (message.event === "NEW_ROUND") {
            this.question.innerHTML = message.question
            this.showQuestion(message.question)
        }

        if (message.event === "WINNER") {
            this.winnerScreen(message)
        }
    }

    winnerScreen(message) {
        this.winner.innerHTML = message.user
        this.answer.innerHTML = message.answer
        this.waitBar.classList.add("hidden")
        this.winnerWrap.classList.remove("hidden")
        this.countdown.stop()
        this.updateLeaderboard(message)

        this.topPage.classList.replace("default-background", "win-background")
        this.timerTitle.innerHTML = "Bom trabalho chat!"
        this.timerMessage.innerHTML = "Aguardando a próxima rodada..."

        const render = (time) => {
            let seconds = time.totalInSeconds
            if (seconds < 10) {
                seconds = this.formatSeconds(seconds)
            }
            if (seconds <= 5) {
                seconds = `<spam class="color-red">${seconds}</spam>`
            }
            this.timer.innerHTML = seconds
        }
        const complete = () => {
            this.startRound()
        }

        this.countdown = new CountDown(this.getRoundTimer(this.breaktimeTime), render, complete)
    }

    looseScreen() {
        this.conn.send(JSON.stringify({command: "ROUND_END"}))

        this.topPage.classList.replace("default-background", "loose-background")
        this.timerTitle.innerHTML = "Ninguém acertou desta vez 🫤"
        this.timerMessage.innerHTML = "Aguardando a próxima rodada..."

        const render = (time) => {
            let seconds = time.totalInSeconds
            if (seconds < 10) {
                seconds = this.formatSeconds(seconds)
            }
            if (seconds <= 5) {
                seconds = `<spam class="color-red">${seconds}</spam>`
            }
            this.timer.innerHTML = seconds
        }
        const complete = () => {
            this.startRound()
        }

        this.countdown = new CountDown(this.getRoundTimer(this.breaktimeTime), render, complete)
    }

    updateLeaderboard(winner) {
        this.updateWinnerList(winner)

        this.leaderboard.innerHTML = ""
        for (const winner in this.winnersList) {
            let node = document.createElement("li")
            node.innerHTML = `<spam class="counter">${this.winnersList[winner].wins}</spam> <spam class="user">${winner}</spam>`
            this.leaderboard.appendChild(node)
        }
    }

    updateWinnerList(winner) {
        if (this.winnersList[winner.user] !== undefined) {
            this.winnersList[winner.user].wins++
        } else {
            this.winnersList[winner.user] = {wins: 1}
        }

        this.winnersList.sort((a, b) => (b.wins > a.wins) ? 1 : -1)
    }

    getRoundTimer(seconds) {
        let date = new Date();
        date.setSeconds(date.getSeconds() + seconds)
        return date
    }

    showQuestion(question) {
        const delay = 350

        let words = question.split(" ")

        this.question.innerHTML = ""
        let target = this.question

        let interval

        setTimeout(() => {
            function nextWord() {
                if (words.length <= 0) {
                    clearInterval(interval)
                    return;
                }
                let word = " " + `<spam class="word">${words.shift()}</spam>`
                var node = document.createElement("spam");

                node.innerHTML = word

                target.appendChild(node)
            }

            nextWord()
            interval = setInterval(nextWord, delay)
        }, delay)
    }

    stop() {
        this.countdown.stop()
    }

    formatSeconds(t) {
        return t < 10 ? '0' + t : t;
    };
}

export default Game
