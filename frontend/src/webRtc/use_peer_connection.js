import { ref, shallowRef } from "vue";

//remotePeer {
//  peerId
//  pc
//  stream
//  state
//}

export function usePeerConnection(signaling, localStream) {
  const peers = shallowRef(new Map()) // [string,remotePeer]
  const error = ref(null)

  const rtcConfig = {
    iceServers: [
      { urls: 'stun:stun.l.google.com:19302' }
    ]
  }

  function createPeer(peerId) {
    //const connection = new RTCPeerConnection(config)
    const pc = new RTCPeerConnection()
    const stream = new MediaStream()
    const peer = {
      peerId,
      pc,
      stream,
      state: 'new'
    }

    pc.ontrack = (e) => {
      const remote = e.streams[0]
      remote.getTracks().forEach(track => {
        stream.addTrack(track)
      })
    }

    pc.onicecandidate = (e) => {
      if (e.candidate) {
        signaling.send({
          type: 'candidate',
          data: {
            candidate: e.candidate,
            peer_id: peerId
          }
        })
      }
    }

    pc.onconnectionstatechange = () => {
      peer.state = pc.connectionState

      if (pc.connectionState === 'failed') {
        restartIce(peer)
      }
    }

    peers.value.set(peerId, peer)

    return peer
  }

  function getOrCreatePeer(peerId) {
    return peers.value.get(peerId) || createPeer(peerId)
  }

  async function addLocalTracks(peer) {
    const stream = localStream()
    if (!stream) return

    stream.getTracks().forEach(t => peer.pc.addTrack(t, stream))
  }

  async function createOffer(peerId) {
    const peer = getOrCreatePeer(peerId)
    await addLocalTracks(peer)

    const offer = await peer.pc.createOffer()
    await peer.pc.setLocalDescription(offer)

    signaling.send({
      type: 'offer',
      data: {
        offer,
        peer_id: peerId
      }
    })
  }

  async function handleOffer(peerId, offer) {
    const peer = getOrCreatePeer(peerId)
    await addLocalTracks(peer)

    await peer.pc.setRemoteDescription(offer)
    const answer = await peer.pc.createAnswer()
    await peer.pc.setLocalDescription(answer)

    signaling.send({
      type: 'answer',
      data: {
        answer,
        peer_id: peerId
      }
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
    } catch {
      // 远端描述未就绪时，生产环境可暂存后重试
    }
  }

  function restartIce(peer) {
    const pc = peer.pc
    if (pc.signalingState !== 'stable') return
    pc.restartIce()
  }

  function removePeer(peerId) {
    const peer = peers.value.get(peerId)
    if (peer) {
      peer.pc.close()
      peers.value.delete(peerId)
    }
  }

  function closeAll() {
    for (const peer of peers.value.values()) {
      peer.pc.close()
    }
    peers.value.clear()
  }

  return {
    peers,
    error,
    createPeer,
    getOrCreatePeer,
    createOffer,
    handleOffer,
    handleAnswer,
    addCandidate,
    removePeer,
    closeAll
  }
}
