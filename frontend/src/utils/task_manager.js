/*
interface Job {
  key(){}
  delay() {}
  execute() {}
  merge(job) {}
  onComplete() {}
}
*/

class DelayTaskMananger {
  delayTask = new Map(); // {key:value}

  submit(job) {
    let o = this.delayTask.get(job.key())
    // 存在合并
    if (o) {
      o.job.merge(job)
      return
    }

    // 执行任务
    let timer = setTimeout(async () => {
      this.delayTask.delete(job.key())
      await job.execute()
      job.onComplete()
    }, job.delay())

    // 加入延迟
    this.delayTask.set(job.key(), { job, timer })
  }
}

export const delayTaskMananger = new DelayTaskMananger()
