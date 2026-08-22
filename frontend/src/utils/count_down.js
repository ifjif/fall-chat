function countDown(date, callBack, endCallback) {
  let target = date.getTime()

  let timer = setInterval(() => {
    let now = new Date().getTime()

    const distance = target - now

    if (distance < 0) {
      clearInterval(timer)
      endCallback()
      return
    }

    const days = Math.floor(distance / (1000 * 60 * 60 * 24))
    const hours = Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
    const minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60))
    const seconds = Math.floor((distance % (1000 * 60)) / 1000)

    callBack(days, hours, minutes, seconds)
  })
}

export {
  countDown
}
