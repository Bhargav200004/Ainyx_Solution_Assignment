import { useState, useEffect } from 'react'

export default function UserForm({ user, onSave, onCancel }) {
  const [name, setName] = useState('')
  const [dob, setDob] = useState('')
  const [errors, setErrors] = useState({})
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    if (user) {
      setName(user.name)
      setDob(user.dob)
    }
  }, [user])

  function validate() {
    const errs = {}
    if (!name.trim()) errs.name = 'Name is required'
    else if (name.trim().length > 255) errs.name = 'Name must be at most 255 characters'
    if (!dob) errs.dob = 'Date of birth is required'
    else if (new Date(dob) > new Date()) errs.dob = 'Date of birth cannot be in the future'
    return errs
  }

  async function handleSubmit(e) {
    e.preventDefault()
    const errs = validate()
    if (Object.keys(errs).length) {
      setErrors(errs)
      return
    }
    setErrors({})
    setSubmitting(true)
    try {
      await onSave(name.trim(), dob)
    } catch (err) {
      if (err.details) {
        setErrors(err.details)
      } else {
        setErrors({ name: err.message })
      }
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <h2>{user ? 'Edit User' : 'Add User'}</h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="user-name">Name</label>
            <input
              id="user-name"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Enter full name"
              autoFocus
            />
            {errors.name && <div className="error-text">{errors.name}</div>}
          </div>
          <div className="form-group">
            <label htmlFor="user-dob">Date of Birth</label>
            <input
              id="user-dob"
              type="date"
              value={dob}
              onChange={(e) => setDob(e.target.value)}
              max={new Date().toISOString().split('T')[0]}
            />
            {errors.dob && <div className="error-text">{errors.dob}</div>}
          </div>
          <div className="modal-actions">
            <button type="button" className="btn btn-ghost" onClick={onCancel} disabled={submitting}>
              Cancel
            </button>
            <button type="submit" className="btn btn-primary" disabled={submitting}>
              {submitting ? 'Saving…' : user ? 'Update' : 'Create'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
