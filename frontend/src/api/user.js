import http from '../request/http'

export const loginApi = (data) => {
  return http.post('/user/login', data)
}

export const forgetPasswordApi = (data) => {
  return http.post('/user/forgetPassword', data)
}

export const registerApi = (data) => {
  return http.post('/user/register', data)
}

export const resetPasswordApi = (data) => {
  return http.post('/user/resetPassword', data)
}

export const updateProfileApi = (data) => {
  return http.post('/user/updateProfile', data)
}

export const userProfileApi = () => {
  let token = localStorage.getItem('token')
  return http.get("/user/profile", {
    token
  })
}
