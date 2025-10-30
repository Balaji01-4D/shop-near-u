import { createContext, useCallback, useContext, useEffect, useState } from 'react'
import type { ReactNode } from 'react'
import { changePassword, deleteAccount, fetchCurrentUser, loginUser, logoutUser, registerUser } from '../services/userApi'
import type { ChangePasswordPayload, CredentialsPayload, RegisterPayload, User } from '../types/user'

export type AuthContextValue = {
  user: User | null
  loading: boolean
  login: (payload: CredentialsPayload) => Promise<User>
  register: (payload: RegisterPayload) => Promise<User>
  logout: () => Promise<void>
  refresh: () => Promise<User | null>
  changePassword: (payload: ChangePasswordPayload) => Promise<string>
  deleteAccount: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  const hydrate = useCallback(async () => {
    try {
      const current = await fetchCurrentUser()
      setUser(current)
      return current
    } catch (error) {
      setUser(null)
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void hydrate()
  }, [hydrate])

  const login = useCallback(async (payload: CredentialsPayload) => {
    const current = await loginUser(payload)
    setUser(current)
    return current
  }, [])

  const register = useCallback(async (payload: RegisterPayload) => {
    const current = await registerUser(payload)
    setUser(current)
    return current
  }, [])

  const logout = useCallback(async () => {
    await logoutUser()
    setUser(null)
  }, [])

  const refresh = useCallback(async () => {
    return hydrate()
  }, [hydrate])

  const handlePasswordChange = useCallback(async (payload: ChangePasswordPayload) => {
    const responseMessage = await changePassword(payload)
    return responseMessage
  }, [])

  const handleDeleteAccount = useCallback(async () => {
    await deleteAccount()
    setUser(null)
  }, [])

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        login,
        register,
        logout,
        refresh,
        changePassword: handlePasswordChange,
        deleteAccount: handleDeleteAccount,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used inside AuthProvider')
  }
  return context
}
