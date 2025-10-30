import { zodResolver } from '@hookform/resolvers/zod'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'
import { z } from 'zod'
import { useAuth } from '../context/AuthContext'

type PasswordFormValues = z.infer<typeof passwordSchema>

const passwordSchema = z
  .object({
    oldPassword: z.string().min(1, 'Enter your current password'),
    newPassword: z.string().min(6, 'Minimum 6 characters'),
    confirmPassword: z.string().min(6, 'Confirm your password'),
  })
  .refine((values) => values.newPassword === values.confirmPassword, {
    path: ['confirmPassword'],
    message: 'Passwords must match',
  })

export function DashboardPage() {
  const { user, changePassword, deleteAccount } = useAuth()
  const navigate = useNavigate()
  const [passwordFeedback, setPasswordFeedback] = useState<{ tone: 'success' | 'error'; message: string } | null>(null)
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false)
  const [deleteError, setDeleteError] = useState<string | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<PasswordFormValues>({
    resolver: zodResolver(passwordSchema),
    defaultValues: {
      oldPassword: '',
      newPassword: '',
      confirmPassword: '',
    },
  })

  const onChangePassword = async (values: PasswordFormValues) => {
    setPasswordFeedback(null)
    try {
      const message = await changePassword({ oldPassword: values.oldPassword, newPassword: values.newPassword })
      setPasswordFeedback({ tone: 'success', message: message || 'Password updated successfully' })
      reset()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to change your password right now'
      setPasswordFeedback({ tone: 'error', message })
    }
  }

  const handleDeleteAccount = async () => {
    setDeleteError(null)
    setIsDeleting(true)
    try {
      await deleteAccount()
      navigate('/register', { replace: true })
    } catch (error) {
      const message = error instanceof Error ? error.message : 'We could not delete your account. Please try again.'
      setDeleteError(message)
    } finally {
      setIsDeleting(false)
      setShowDeleteConfirm(false)
    }
  }

  return (
    <div className="grid-two">
      <section className="card" style={{ gap: '1.75rem', display: 'flex', flexDirection: 'column' }}>
        <div>
          <div className="hero-eyebrow">Your profile</div>
          <h2 style={{ marginBottom: '0.5rem' }}>Welcome back, {user?.name || user?.email}</h2>
          <p className="text-muted">Manage account security and keep your details up to date.</p>
        </div>

        <div className="subtle-card" style={{ display: 'grid', gap: '0.6rem' }}>
          <div>
            <strong>Email</strong>
            <p className="text-muted" style={{ marginTop: '0.25rem' }}>{user?.email}</p>
          </div>
          <div>
            <strong>User ID</strong>
            <p className="text-muted" style={{ marginTop: '0.25rem' }}>#{user?.id}</p>
          </div>
        </div>

        <div>
          <h3 className="section-title" style={{ fontSize: '1.1rem' }}>Change password</h3>
          <p className="section-subtitle" style={{ marginBottom: '1.5rem' }}>
            Keep your account protected with a strong password you refresh regularly.
          </p>

          <form onSubmit={handleSubmit(onChangePassword)} className="form-grid" noValidate>
            <div>
              <label htmlFor="oldPassword">Current password</label>
              <input id="oldPassword" type="password" autoComplete="current-password" {...register('oldPassword')} />
              {errors.oldPassword ? <p className="error-text">{errors.oldPassword.message}</p> : null}
            </div>

            <div>
              <label htmlFor="newPassword">New password</label>
              <input id="newPassword" type="password" autoComplete="new-password" {...register('newPassword')} />
              {errors.newPassword ? <p className="error-text">{errors.newPassword.message}</p> : null}
            </div>

            <div>
              <label htmlFor="confirmPassword">Confirm new password</label>
              <input id="confirmPassword" type="password" autoComplete="new-password" {...register('confirmPassword')} />
              {errors.confirmPassword ? <p className="error-text">{errors.confirmPassword.message}</p> : null}
            </div>

            {passwordFeedback ? (
              <div className={`alert ${passwordFeedback.tone === 'success' ? 'alert-success' : ''}`}>
                {passwordFeedback.message}
              </div>
            ) : null}

            <button type="submit" className="btn btn-primary" disabled={isSubmitting}>
              {isSubmitting ? 'Updating…' : 'Update password'}
            </button>
          </form>
        </div>
      </section>

      <section className="card" style={{ display: 'flex', flexDirection: 'column', gap: '1.75rem' }}>
        <div>
          <h2 style={{ marginBottom: '0.75rem' }}>Account controls</h2>
          <p className="text-muted">
            Prefer not to continue? You are always in control of your data and access.
          </p>
        </div>

        <div className="subtle-card" style={{ display: 'grid', gap: '0.9rem' }}>
          <div>
            <strong>Need a break?</strong>
            <p className="text-muted" style={{ marginTop: '0.35rem' }}>
              You can sign out anytime from the header. Your personalised recommendations stay ready when you return.
            </p>
          </div>
        </div>

        <div className="subtle-card" style={{ border: '1px solid rgba(248, 113, 113, 0.35)', background: 'rgba(254, 226, 226, 0.45)' }}>
          <strong>Delete account</strong>
          <p className="text-muted" style={{ margin: '0.5rem 0 1rem' }}>
            This will permanently remove your profile and sign you out on all devices.
          </p>

          {deleteError ? <div className="alert">{deleteError}</div> : null}

          {showDeleteConfirm ? (
            <div style={{ display: 'grid', gap: '0.75rem' }}>
              <p style={{ color: '#b91c1c', margin: 0 }}>Are you sure? This action cannot be reversed.</p>
              <div style={{ display: 'flex', gap: '0.75rem', flexWrap: 'wrap' }}>
                <button type="button" className="btn btn-danger" onClick={handleDeleteAccount} disabled={isDeleting}>
                  {isDeleting ? 'Deleting…' : 'Yes, delete my account'}
                </button>
                <button type="button" className="btn btn-ghost" onClick={() => setShowDeleteConfirm(false)} disabled={isDeleting}>
                  Keep my account
                </button>
              </div>
            </div>
          ) : (
            <button type="button" className="btn btn-danger" onClick={() => setShowDeleteConfirm(true)}>
              Delete account
            </button>
          )}
        </div>
      </section>
    </div>
  )
}
