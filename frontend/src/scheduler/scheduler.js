export class Scheduler {
  maxConcurrent;
  runningCount;
  queue = [];
  isPaused;

  constructor(maxConcurrent = 3) {
    this.maxConcurrent = maxConcurrent
  }

  add(run) {
    return new Promise((resolve, reject) => {
      this.queue.push({ run, resolve, reject })
      this.runNext()
    })
  }

  paused() {
    this.isPaused = true
  }

  resume() {
    this.isPaused = false
    this.runNext()
  }

  cancel() {
    this.queue = []
    this.runningCount = 0
    this.isPaused = false
  }

  getQueueLength() {
    return this.queue.length
  }

  runNext() {
    if (this.isPaused || this.runningCount >= this.maxConcurrent || this.queue.length === 0) return

    const task = this.queue.shift()
    this.runningCount++

    task.run()
      .then(task.resolve)
      .catch(task.reject)
      .finally(() => {
        this.runningCount--
        this.runNext()
      })
  }
}
