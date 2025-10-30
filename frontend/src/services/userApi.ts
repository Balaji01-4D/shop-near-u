import { apiRequest } from './apiClient'
import type { ChangePasswordPayload, CredentialsPayload, RegisterPayload, User } from '../types/user'

type AuthSuccess = {
  user: User
  token: string
}

export async function loginUser(payload: CredentialsPayload) {
  const { data } = await apiRequest<AuthSuccess>('/auth/login', {
    method: 'POST',
    body: {
      email: payload.email,
      password: payload.password,
    },
  })

  return data.user
}

export async function registerUser(payload: RegisterPayload) {
  const { data } = await apiRequest<AuthSuccess>('/auth/register', {
    method: 'POST',
    body: {
      name: payload.name,
      email: payload.email,
      password: payload.password,
    },
  })

  return data.user
}

export async function fetchCurrentUser() {
  const { data } = await apiRequest<User>('/auth/me', {
    method: 'GET',
  })
  return data
}

export async function logoutUser() {
  await apiRequest<null>('/auth/logout', {
    method: 'POST',
  })
}

export async function changePassword(payload: ChangePasswordPayload) {
  const { message } = await apiRequest<null>('/auth/change-password', {
    method: 'POST',
    body: {
      old_password: payload.oldPassword,
      new_password: payload.newPassword,
    },
  })

  return message
}

export async function deleteAccount() {
  await apiRequest<null>('/auth/delete-account', {
    method: 'DELETE',
  })
}
