import { zodResolver } from '@hookform/resolvers/zod'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { Link, useNavigate } from 'react-router-dom'
import { z } from 'zod'
import { useAuth } from '../context/AuthContext'

type FormValues = z.infer<typeof schema>

const schema = z
  .object({
    name: z.string().min(2, 'Tell us your name'),
    email: z.string().email('Enter a valid email'),
    password: z.string().min(6, 'At least 6 characters for security'),
    confirmPassword: z.string().min(6, 'Please confirm your password'),
  })
  .refine((values) => values.password === values.confirmPassword, {
    message: 'Passwords need to match',
    path: ['confirmPassword'],
  })

export function RegisterPage() {
  const { register: registerUser } = useAuth()
  const navigate = useNavigate()
  const [serverError, setServerError] = useState<string | null>(null)

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: '',
      email: '',
      password: '',
      confirmPassword: '',
    },
  })

  const onSubmit = async (values: FormValues) => {
    setServerError(null)
    try {
      await registerUser({
        name: values.name,
        email: values.email,
        password: values.password,
      })
      reset()
      navigate('/dashboard', { replace: true })
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to create your account right now'
      setServerError(message)
    }
  }

  return (
    <div className="card">
      <div className="pill" style={{ width: 'fit-content', marginBottom: '1.25rem' }}>
        Start exploring
      </div>
      <h2>Create your account</h2>
      <p className="text-muted">
        Unlock the full ShopNearU experience and never miss out on neighbourhood discoveries.
      </p>

      <form onSubmit={handleSubmit(onSubmit)} className="form-grid" noValidate>
        <div>
          <label htmlFor="name">Full name</label>
          <input id="name" type="text" autoComplete="name" {...register('name')} />
          {errors.name ? <p className="error-text">{errors.name.message}</p> : null}
        </div>

        <div>
          <label htmlFor="email">Email</label>
          <input id="email" type="email" autoComplete="email" {...register('email')} />
          {errors.email ? <p className="error-text">{errors.email.message}</p> : null}
        </div>

        <div>
          <label htmlFor="password">Password</label>
          <input id="password" type="password" autoComplete="new-password" {...register('password')} />
          {errors.password ? <p className="error-text">{errors.password.message}</p> : null}
        </div>

        <div>
          <label htmlFor="confirmPassword">Confirm password</label>
          <input id="confirmPassword" type="password" autoComplete="new-password" {...register('confirmPassword')} />
          {errors.confirmPassword ? <p className="error-text">{errors.confirmPassword.message}</p> : null}
        </div>

        {serverError ? <div className="alert">{serverError}</div> : null}

        <button type="submit" className="btn btn-primary" disabled={isSubmitting}>
          {isSubmitting ? 'Creating account…' : 'Create account'}
        </button>
      </form>

      <p className="text-muted" style={{ marginTop: '1.75rem' }}>
        Already part of the community?{' '}
        <Link className="muted-link" to="/login">
          Sign in
        </Link>
      </p>
    </div>
  )
}
