import { fileAccessUrlApi } from "@/api/file"

export const getRemoteFileUrl = (uri) => {
  return fileAccessUrlApi({ file: uri })
}
