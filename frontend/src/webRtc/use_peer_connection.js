import { ref, shallowRef } from 'vue'
// 维护对端连接
// remotePeer {
//  pc
//  peerId
//  stream
//  state
// }
export function usePeerConnection(signaling, localStream) {
  const peers = shallowRef(new Map())// {peerId:remotePeer}
  const error = ref(null)

  const rtcConfig = {
    iceServers: [
      //{ urls: "" }  // 配置 stun 和 turn
    ]
  }

  function createPeer(peerId) {
    const pc = new RTCPeerConnection(rtcConfig)
    const stream = new MediaStream()
    const peer = {
      peerId,
      pc,
      stream,
      state: 'new'
    }

    // 有track 后加入 stream
    pc.ontrack = (e) => {
      const remote = e.streams[0]
      remote.getTracks().forEach(t => {
        stream.addTrack(t)
      })
    }

    // ice交换, 本端有新的了，发送给对端
    pc.onicecandidate = (e) => {
      if (e.candidate) {
        signaling.sendCandidate({
          candidate: e.candidate,
          peer_id: peerId
        })
      }
    }

    // 连接状态改变
    pc.onconnectionstatechange = () => {
      peer.state = pc.connectionState
      if (peer.state === 'failed') {
        restartIce(peer)
      }
    }

    peers.value.set(peerId, peer)
    return peer
  }

  function restartIce(peer) {
    const pc = peer.pc
    if (pc.signalingState !== 'stable') return
    pc.restartIce()
  }

  function removePeer(peerId) {
    const peer = peers.value.get(peerId)
    if (!peer) return
    peer.pc.close()
    peers.value.delete(peerId)
  }

  function closeAll() {
    for (const peer of peers.value.values()) {
      peer.pc.close()
    }
    peers.value.clear()
  }

  // 避免重复添加
  function addLocalTracks(peer) {
    const stream = localStream()
    if (!stream) return
    stream.getTracks().forEach(t => {
      const exists = peer.pc.getSenders().some(s => s.track && s.track.kind === t.kind)
      if (!exists) peer.pc.addTrack(t, stream)
    })
  }

  // 将一条track 给所有 peer替换
  function replaceLocalTrack(track) {
    if (!track) return
    for (const peer of peers.value.values()) {
      const sender = peer.pc.getSenders().find(s => s.track && s.track.kind === track.kind)
      if (sender) {
        sender.replaceTrack(track)
      }
    }
  }

  function getOrCreatePeer(peerId) {
    return peers.value.get(peerId) || createPeer(peerId)
  }

  async function createOffer(peerId) {
    const peer = getOrCreatePeer(peerId)
    addLocalTracks(peer)
    const offer = await peer.pc.createOffer()
    await peer.pc.setLocalDescription(offer)
    signaling.sendOffer({
      peer_id: peerId,
      offer
    })

  }
  async function handleOffer(peerId, offer) {
    const peer = getOrCreatePeer(peerId)
    addLocalTracks(peer)

    await peer.pc.setRemoteDescription(offer)
    const answer = await peer.pc.createAnswer()
    await peer.pc.setLocalDescription(answer)

    signaling.sendAnswer({
      peer_id: peerId,
      answer
    })
  }
  async function handleAnswer(peerId, answer) {
    const peer = peers.value.get(peerId)
    if (!peer) return
    await peer.pc.setRemoteDescription(answer)
  }

  async function addCandidate(peerId, candidate) {
    const peer = peers.value.get(peerId)
    if (!peer) return

    try {
      await peer.pc.addIceCandidate(candidate)
    } catch (e) {
      // 远端描述未就绪时，生产环境可暂存后重试
      console.log('[add candidate]', e)
    }
  }



  return {
    peers,
    error,
    removePeer,
    closeAll,
    createPeer,
    createOffer,
    handleOffer,
    handleAnswer,
    addCandidate,
    replaceLocalTrack
  }
}
