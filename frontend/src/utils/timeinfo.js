import { padZero } from "./pad_zero"
// 解析出聊天时间
// 1) 1分钟内 -> 显示 刚刚
// 2) 1小时内 -> 显示 xxx分钟前
// 3) 1天内 -> 显示 (上午/中午/下午/晚上 xx:xx)
// 4) 1周内 -> 显示(周几 上午/中午/下午/晚上 xx:xx)
// 5) 1年内 -> 显示(x月x日 上午/中午/下午/晚上 xx:xx)

// 6) 其它  -> 显示(x年x月x日 上午/中午/下午/晚上 xx:xx)
export function parseTime(t) {
  if (!t) t = Date.now()
  const now = new Date()
  const date = new Date(t)
  const diff = now - date // 毫秒差

  //获取时段前缀
  function getPeriod(h) {
    if (h >= 5 && h < 12) return '上午'
    if (h >= 12 && h < 14) return '中午'
    if (h >= 14 && h < 18) return '下午'
    return '晚上'
  }

  //1) 1分钟内
  if (diff < 60 * 1000) {
    return '刚刚'
  }

  // 2) 1小时内
  if (diff < 24 * 60 * 1000) {
    return Math.floor(diff / (60 * 1000)) + "分钟前"
  }

  const h = padZero(date.getHours())
  const m = padZero(date.getMinutes())
  const period = getPeriod(h)
  const timeStr = `${period} ${h}:${m}`

  // 3) 1天内
  if (
    now.getFullYear() === date.getFullYear() &&
    now.getMonth() === date.getMonth() &&
    now.getDate() === date.getDate()
  ) {
    return timeStr
  }

  // 4) 1周内
  const oneWeek = 7 * 24 * 60 * 60 * 1000
  if (diff < oneWeek) {
    const weekDays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
    return `${weekDays[date.getDay()]} ${timeStr}`
  }

  // 5) 1年内
  if (now.getFullYear() === date.getFullYear()) {
    return `${date.getMonth() + 1}月${date.getDate()}日 ${timeStr}`
  }

  // 6) 其它
  return `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日 ${timeStr}`
}
