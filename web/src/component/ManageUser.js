import { useState, useEffect } from "react";
import { Link } from "react-router-dom";

function ManageUser() {
  const [users, setUsers] = useState([]);
  const roleDict = {
    0: "Anonymous",
    1: "Admin",
    2: "Normal User",
  };

  function getUsers() {
    fetch("/api/v1/get_users")
      .then((res) => res.json())
      .then((data) => {
        if (data.error) {
          alert(data.error);
        } else {
          setUsers(data.users);
        }
      });
  }

  useEffect(() => {
    getUsers();
  }, []);

  return (
    <div className="page">
      <h3>Manage User</h3>
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Role</th>
            <th>Active</th>
          </tr>
        </thead>
        <tbody>
          {users.map((user) => (
            <tr key={user.id}>
              <td>
                <Link to={`/manage/users/${user.id}`}>@{user.username}</Link>
              </td>
              <td>{roleDict[user.role]}</td>
              <td>
                <input
                  type="checkbox"
                  defaultChecked={user.active}
                  onClick={(e) => {
                    fetch("/api/v1/update_user_active", {
                      method: "POST",
                      headers: {
                        "Content-Type": "application/json",
                      },
                      body: JSON.stringify({
                        id: user.id,
                        active: e.target.checked,
                      }),
                    }).then((res) => res.json()).then((data) => {
                      if (data.error) {
                        alert(data.error);
                      } else {
                        getUsers();
                      }
                    });
                  }}
                />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default ManageUser;
