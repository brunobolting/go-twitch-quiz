class CountDown {
    constructor(expiredDate, onRender, onComplete) {
        this.setExpiredDate(expiredDate)

        this.onRender = onRender
        this.onComplete = onComplete
    }

    setExpiredDate(expiredDate) {
        // get the current time
        const currentTime = new Date().getTime();

        // calculate the remaining time
        this.timeRemaining = expiredDate.getTime() - currentTime;

        // should the countdown completes or start
        this.timeRemaining > 0 ?
            this.start() :
            this.complete();
    }

    complete() {
        if (typeof this.onComplete === 'function') {
            this.onComplete();
        }
    }

    start() {
        // update the countdown
        this.update();

        //  setup a timer
        this.intervalId = setInterval(() => {
            // update the timer
            this.timeRemaining -= 1000;

            if (this.timeRemaining < 0) {
                // call the callback
                this.complete();

                // clear the interval if expired
                clearInterval(this.intervalId);
            } else {
                this.update();
            }

        }, 1000);
    }

    getTime() {
        return {
            days: Math.floor(this.timeRemaining / 1000 / 60 / 60 / 24),
            hours: Math.floor(this.timeRemaining / 1000 / 60 / 60) % 24,
            minutes: Math.floor(this.timeRemaining / 1000 / 60) % 60,
            seconds: Math.floor(this.timeRemaining / 1000) % 60,
            totalInSeconds: Math.floor(this.timeRemaining / 1000)
        };
    }

    update() {
        if (typeof this.onRender === 'function') {
            this.onRender(this.getTime());
        }
    }

    stop() {
        clearInterval(this.intervalId);
    }
}

export default CountDown
