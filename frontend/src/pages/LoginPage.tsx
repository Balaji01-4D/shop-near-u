import { zodResolver } from '@hookform/resolvers/zod'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { z } from 'zod'
import { useAuth } from '../context/AuthContext'

type FormValues = z.infer<typeof schema>

const schema = z.object({
  email: z.string().min(1, 'Email is required').email('Enter a valid email'),
  password: z.string().min(1, 'Password is required'),
})

export function LoginPage() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const [serverError, setServerError] = useState<string | null>(null)

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      email: '',
      password: '',
    },
  })

  const onSubmit = async (values: FormValues) => {
    setServerError(null)
    try {
      await login({ email: values.email, password: values.password })
      const redirectTo = (location.state as { from?: { pathname?: string } } | null)?.from?.pathname ?? '/dashboard'
      navigate(redirectTo, { replace: true })
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to sign in right now'
      setServerError(message)
    }
  }

  return (
    <div className="card">
      <div className="pill" style={{ width: 'fit-content', marginBottom: '1.25rem' }}>
        Welcome back
      </div>
      <h2>Sign in to ShopNearU</h2>
      <p className="text-muted">
        Access curated shops, manage your subscriptions, and keep track of your favourite stores.
      </p>

      <form onSubmit={handleSubmit(onSubmit)} className="form-grid" noValidate>
        <div>
          <label htmlFor="email">Email</label>
          <input id="email" type="email" autoComplete="email" {...register('email')} />
          {errors.email ? <p className="error-text">{errors.email.message}</p> : null}
        </div>

        <div>
          <label htmlFor="password">Password</label>
          <input id="password" type="password" autoComplete="current-password" {...register('password')} />
          {errors.password ? <p className="error-text">{errors.password.message}</p> : null}
        </div>

        {serverError ? <div className="alert">{serverError}</div> : null}

        <button type="submit" className="btn btn-primary" disabled={isSubmitting}>
          {isSubmitting ? 'Signing in…' : 'Sign in'}
        </button>
      </form>

      <p className="text-muted" style={{ marginTop: '1.75rem' }}>
        New to ShopNearU?{' '}
        <Link className="muted-link" to="/register">
          Create an account
        </Link>
      </p>
    </div>
  )
}
