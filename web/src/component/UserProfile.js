import { useState, useEffect } from "react";
import { useParams } from "react-router";
import ReviewEntry from "./ReviewEntry";

function UserProfile(props) {
  let params = useParams();
  const [reviews, setReviews] = useState([]);
  const [oldPassword, setOldPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [newPasswordConfirm, setNewPasswordConfirm] = useState("");
  const [user, setUser] = useState({
    id: 0,
    username: "",
    role: 0,
    active: false,
    avatar_id: 0,
  });

  function getReviews() {
    fetch("/api/v1/get_reviews_by_user", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        id: parseInt(params.id),
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          alert(data.error);
        } else {
          setReviews(data);
        }
      });
  }

  function getUserInfo() {
    fetch("/api/v1/get_user_info", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        id: parseInt(params.id),
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          alert(data.error);
        } else {
          setUser(data.user);
        }
      });
  }

  useEffect(() => {
    getReviews();
    getUserInfo();
  }, []);

  return (
    <div className="page">
      <h3>User Profile</h3>
      <div className="horizontal">
        <input
          type="text"
          value={user.username}
          onChange={(e) => {
            setUser({
              ...user,
              username: e.target.value,
            });
          }}
        />
        <button
          onClick={() => {
            fetch("/api/v1/update_username", {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
              },
              body: JSON.stringify({
                id: parseInt(params.id),
                username: user.username,
              }),
            })
              .then((response) => response.json())
              .then((data) => {
                if (data.error) {
                  alert(data.error);
                } else {
                  props.setUser({
                    ...props.user,
                    username: user.username,
                  });
                  alert("Username updated successfully!");
                  getUserInfo();
                }
              });
          }}
          disabled={props.user.id !== user.id && props.user.role !== 1}
        >
          Save Username
        </button>
      </div>
      <div>
        <input
          type="password"
          value={oldPassword}
          placeholder="Old Password"
          onChange={(e) => setOldPassword(e.target.value)}
        />
        <input
          type="password"
          value={newPassword}
          placeholder="New Password"
          onChange={(e) => setNewPassword(e.target.value)}
        />
        <input
          type="password"
          value={newPasswordConfirm}
          placeholder="Confirm New Password"
          onChange={(e) => setNewPasswordConfirm(e.target.value)}
        />
        <button
          onClick={() => {
            fetch("/api/v1/update_user_password", {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
              },
              body: JSON.stringify({
                id: parseInt(params.id),
                old_password: oldPassword,
                new_password: newPassword,
                new_password_confirm: newPasswordConfirm,
              }),
            })
              .then((response) => response.json())
              .then((data) => {
                if (data.error) {
                  alert(data.error);
                } else {
                  alert("Password updated successfully!");
                }
              });
          }}
          disabled={
            (props.user.id !== user.id && props.user.role !== 1) ||
            newPassword !== newPasswordConfirm ||
            newPassword.length === 0
          }
        >
          Change Password
        </button>
      </div>
      <h4>Reviews</h4>
      {reviews.map((review) => (
        <ReviewEntry key={review.id} review={review} user={props.user} />
      ))}
    </div>
  );
}

export default UserProfile;
