export default function UserTable({ users, onEdit, onDelete }) {
  if (!users || users.length === 0) {
    return <div className="empty-state">No users found. Add one to get started.</div>
  }

  return (
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Name</th>
          <th>Date of Birth</th>
          <th>Age</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        {users.map((user) => (
          <tr key={user.id}>
            <td>{user.id}</td>
            <td>{user.name}</td>
            <td>{user.dob}</td>
            <td>{user.age ?? '—'}</td>
            <td>
              <div className="actions">
                <button
                  className="btn btn-ghost btn-sm"
                  onClick={() => onEdit(user)}
                  title="Edit user"
                  id={`edit-user-${user.id}`}
                >
                  ✏️
                </button>
                <button
                  className="btn btn-danger btn-sm"
                  onClick={() => onDelete(user)}
                  title="Delete user"
                  id={`delete-user-${user.id}`}
                >
                  🗑️
                </button>
              </div>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
