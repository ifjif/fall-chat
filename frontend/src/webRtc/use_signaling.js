// 信令服务
// {
//  type: offer | answer | candidate | join | leave
//  data{
//  offer | answer | candidate
//  peerId:
//  roomId:
//  data:
//  }
// }
import { useWsStore } from "@/stores/ws_store";
export function useSignaling() {
  const wsStore = useWsStore()

  function send(msg) {
    wsStore.sendMessage(msg)
  }
  function on(type, handler) {
    wsStore.onOnce(type, msg => {
      handler?.(msg.data)
    })
  }

  // {peer_id, candidate}
  function sendOffer(data) {
    send({
      type: 'offer',
      data
    })
  }

  // {peer_id, data}
  function onOffer(handler) {
    on('offer', handler)
  }

  // {peer_id, answer}
  function sendAnswer(data) {
    send({
      type: 'answer',
      data
    })
  }
  // {peer_id, data}
  function onAnswer(handler) {
    on('answer', handler)
  }

  // {peer_id, candidate}
  function sendCandidate(data) {
    send({
      type: 'candidate',
      data
    })
  }

  // {peer_id, data}
  function onCandidate(handler) {
    on('candidate', handler)
  }

  return {
    send,
    on,
    sendOffer,
    onOffer,
    sendAnswer,
    onAnswer,
    sendCandidate,
    onCandidate
  }
}
