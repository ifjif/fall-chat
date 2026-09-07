// websocket 全局工具
// 状态：
//   closed | connected | error
// 回调：
//    (msg) => void

import { wsAccessUrlApi } from "@/api/wsc";

export class WebSocketService {
  static instance = null;
  url;
  ws;
  isConnected = false;
  reconnectIntervalTime = 3000;
  maxReconnectAttempts = 10;
  reconnectAttempts = 0;
  heartbeatInterval;
  heartbeatIntervalTime = 5000;
  messageCallbacks = {}; // {key: []}

  constructor() {
    this.connect()
  }

  static getInstance() {
    if (!WebSocketService.instance) WebSocketService.instance = new WebSocketService()
    return WebSocketService.instance
  }

  async connect() {
    if (this.ws) return

    this.url = await wsAccessUrlApi()
    this.ws = new WebSocket(this.url)

    this.ws.onopen = () => {
      console.log('[WebSocket] Connected')
      this.isConnected = true
      this.reconnectAttempts = 0
      this.startHeartbeat()
    }

    this.ws.onclose = (event) => {
      console.log('[WebSocket] Disconnected')
      this.ws = null
      this.isConnected = false
      this.stopHeartbeat()

      // 非正常关闭进行重连
      if (event.code !== 1000 || event.code !== 1001) {
        this.scheduleReconnect()
      }
    }

    this.ws.onmessage = (event) => {
      let data
      try {
        data = JSON.parse(event.data)
      } catch (e) {
        data = { type: 'raw', data: event.data }
      }

      this.dispatchMessage(data)
    }

    this.ws.onerror = (err) => {
      console.log('[WebSocket] Error:', err)
    }
  }

  isReady() {
    if (!this.isConnected) return false
    return this.ws?.readyState === WebSocket.OPEN
  }

  send(data) {
    if (!this.isReady()) return
    this.ws.send(JSON.stringify(data))
  }

  startHeartbeat() {
    this.stopHeartbeat()

    this.heartbeatInterval = setInterval(() => {
      if (this.isReady()) this.send({ type: 'ping' })
    }, this.heartbeatIntervalTime)
  }

  stopHeartbeat() {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval)
    }
  }

  scheduleReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) return
    setTimeout(() => {
      this.reconnectAttempts++
      this.connect()
    }, this.reconnectIntervalTime)
  }

  on(type, callback) {
    if (!this.messageCallbacks[type]) this.messageCallbacks[type] = []
    this.messageCallbacks[type].push(callback)
  }

  off(type, callback) {
    if (!this.messageCallbacks[type]) return
    this.messageCallbacks[type] = this.messageCallbacks[type].filter(cb => cb !== callback)
    if (this.messageCallbacks[type].length === 0) delete (this.messageCallbacks[type])
  }

  dispatchMessage(msg) {
    this.messageCallbacks[msg.type]?.forEach(cb => cb(msg))
    // 全局监听
    this.messageCallbacks['*']?.forEach(cb => cb(msg))
  }
}
