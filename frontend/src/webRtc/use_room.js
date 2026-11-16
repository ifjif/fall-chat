import { ref } from 'vue'
export function useRoom() {
  const selfId = ref('')
  const roomId = ref('')
  const remotePeers = ref([])

  function create(rId, uid) {
    roomId.value = rId
    selfId.value = uid
    remotePeers.value = []
  }

  function addRemote(peerId) {
    if (remotePeers.value.includes(peerId)) return
    remotePeers.value.push(peerId)
  }
  function removeRemote(peerId) {
    remotePeers.value = remotePeers.value.filter(id => id !== peerId)
  }

  function leave() {
    roomId.value = ''
    selfId.value = ''
    remotePeers.value = []
  }

  return {
    selfId,
    roomId,
    remotePeers,
    create,
    addRemote,
    removeRemote,
    leave
  }
}
