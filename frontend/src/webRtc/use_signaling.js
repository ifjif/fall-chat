//export type SignalingMessage = {
//  type: 'offer' | 'answer' | 'candidate' | 'join' | 'leave'
//  roomId?: string
//  peerId?: string
//  data?: any
//}

import { useWsStore } from "@/stores/ws_store";

export function useSignaling() {
  const wsStore = useWsStore()

  function send(msg) {
    wsStore.sendMessage(msg)
  }

  function on(type, handler) {
    wsStore.onOnce(type, msg => {
      handler(msg.data)
    })
  }

  return {
    send,
    on
  }
}
