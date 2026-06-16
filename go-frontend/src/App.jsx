import { useState, useEffect, useCallback } from 'react'
import { listUsers, createUser, updateUser, deleteUser } from './api'
import UserTable from './components/UserTable'
import UserForm from './components/UserForm'
import Toast from './components/Toast'

export default function App() {
  const [users, setUsers] = useState([])
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [toasts, setToasts] = useState([])
  const LIMIT = 10

  // Modal state
  const [showForm, setShowForm] = useState(false)
  const [editingUser, setEditingUser] = useState(null)
  const [confirmDelete, setConfirmDelete] = useState(null)

  const addToast = useCallback((message, type = 'success') => {
    const id = Date.now()
    setToasts((prev) => [...prev, { id, message, type }])
  }, [])

  const removeToast = useCallback((id) => {
    setToasts((prev) => prev.filter((t) => t.id !== id))
  }, [])

  const fetchUsers = useCallback(async () => {
    setLoading(true)
    try {
      const data = await listUsers(page, LIMIT)
      setUsers(data.data || [])
      setTotalPages(data.total_pages || 1)
      setTotal(data.total || 0)
    } catch (err) {
      addToast(err.message, 'error')
    } finally {
      setLoading(false)
    }
  }, [page, addToast])

  useEffect(() => {
    fetchUsers()
  }, [fetchUsers])

  // Create or update
  async function handleSave(name, dob) {
    if (editingUser) {
      await updateUser(editingUser.id, name, dob)
      addToast('User updated successfully')
    } else {
      await createUser(name, dob)
      addToast('User created successfully')
    }
    setShowForm(false)
    setEditingUser(null)
    fetchUsers()
  }

  function handleEdit(user) {
    setEditingUser(user)
    setShowForm(true)
  }

  function handleDeleteClick(user) {
    setConfirmDelete(user)
  }

  async function handleConfirmDelete() {
    try {
      await deleteUser(confirmDelete.id)
      addToast('User deleted successfully')
      setConfirmDelete(null)
      // If we deleted the last item on this page, go back a page
      if (users.length === 1 && page > 1) {
        setPage(page - 1)
      } else {
        fetchUsers()
      }
    } catch (err) {
      addToast(err.message, 'error')
      setConfirmDelete(null)
    }
  }

  return (
    <div className="app">
      <header className="app-header">
        <h1>
          <span>Ainyx</span> Users
        </h1>
        <button
          id="add-user-btn"
          className="btn btn-primary"
          onClick={() => {
            setEditingUser(null)
            setShowForm(true)
          }}
        >
          + Add User
        </button>
      </header>

      <div className="table-wrap">
        {loading ? (
          <div className="loading">
            <div className="spinner" />
          </div>
        ) : (
          <>
            <UserTable users={users} onEdit={handleEdit} onDelete={handleDeleteClick} />
            {total > 0 && (
              <div className="pagination">
                <span>
                  Page {page} of {totalPages} · {total} user{total !== 1 ? 's' : ''}
                </span>
                <div className="pagination-btns">
                  <button
                    className="btn btn-ghost btn-sm"
                    disabled={page <= 1}
                    onClick={() => setPage((p) => p - 1)}
                  >
                    ← Prev
                  </button>
                  <button
                    className="btn btn-ghost btn-sm"
                    disabled={page >= totalPages}
                    onClick={() => setPage((p) => p + 1)}
                  >
                    Next →
                  </button>
                </div>
              </div>
            )}
          </>
        )}
      </div>

      {/* Add / Edit Modal */}
      {showForm && (
        <UserForm
          user={editingUser}
          onSave={handleSave}
          onCancel={() => {
            setShowForm(false)
            setEditingUser(null)
          }}
        />
      )}

      {/* Delete Confirmation Modal */}
      {confirmDelete && (
        <div className="modal-overlay" onClick={() => setConfirmDelete(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h2>Delete User</h2>
            <p className="confirm-text">
              Are you sure you want to delete <strong>{confirmDelete.name}</strong>? This action cannot be undone.
            </p>
            <div className="modal-actions">
              <button className="btn btn-ghost" onClick={() => setConfirmDelete(null)}>
                Cancel
              </button>
              <button className="btn btn-primary" style={{ background: 'var(--danger)' }} onClick={handleConfirmDelete}>
                Delete
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Toasts */}
      <div className="toast-container">
        {toasts.map((t) => (
          <Toast key={t.id} message={t.message} type={t.type} onClose={() => removeToast(t.id)} />
        ))}
      </div>
    </div>
  )
}
