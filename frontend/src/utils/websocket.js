/*
 * 
interface WsMessage {
  type: string;
  [key: string]: any;
}

type MessageCallback = (data: WsMessage) => void;
*
 *
 */

import { wsAccessUrlApi } from "@/api/wsc";

export class WebSocketService2 {
  static instance;
  ws;
  url;
  reconnectInterval = 3000;
  maxReconnectAttempts = 10;
  reconnectAttempts = 0;
  isConnected = false;
  messageCallbacks = {};
  heartbeatInterval;
  heartbeatTimeout = 5000;

  constructor() {
    this.connect()
  }

  static getInstance() {
    if (!WebSocketService2.instance) {
      WebSocketService2.instance = new WebSocketService2()
    }

    return WebSocketService2.instance
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

    this.ws.onmessage = (event) => {
      let data
      try {
        data = JSON.parse(event.data)
      } catch (err) {
        data = { type: 'raw', data: event.data }
      }

      this.dispatchMessage(data)
    }

    this.ws.onclose = (event) => {
      console.log("[WebSocket] DisConnected", event)
      this.isConnected = false
      this.stopHeartbeat()
      this.ws = null
      // 非正常关闭，尝试重连
      if (event.code !== 1000 && event.code !== 1001) {
        this.scheduleReconnect()
      }
    }

    this.ws.onerror = (error) => {
      console.log('[WebSocket] Error:', error)
    }
  }

  send(data) {
    if (this.isConnected && this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data))
    }
  }

  // 订阅特定类型的消息
  on(type, callback) {
    if (!this.messageCallbacks[type]) this.messageCallbacks[type] = []
    this.messageCallbacks[type].push(callback)
  }

  // 取消订阅
  off(type, callback) {
    if (this.messageCallbacks[type]) {
      this.messageCallbacks[type] = this.messageCallbacks[type].filter(cb => cb !== callback)
    }
  }

  dispatchMessage(data) {
    // 触发特定类型的监听器
    this.messageCallbacks[data.type]?.forEach(cb => cb(data))
    // 触发全局通配符（供pinia使用）
    this.messageCallbacks['*']?.forEach(cb => cb(data))
  }

  startHeartbeat() {
    this.stopHeartbeat()
    this.heartbeatInterval = setInterval(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify({ type: 'ping' }))
      }
    }, this.heartbeatTimeout)
  }

  stopHeartbeat() {
    if (this.heartbeatInterval) clearInterval(this.heartbeatInterval)
  }

  scheduleReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) return
    setTimeout(() => {
      this.reconnectAttempts++
      this.connect()
    }, this.reconnectInterval)
  }
}
