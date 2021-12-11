import { useEffect } from 'react';

function UserStatus(props) {
  // props.user
  // props.setUser
  useEffect(() => {
    fetch("/api/v1/login")
      .then(res => res.json())
      .then(data => {
        props.setUser(data.user);
      });
  }, []);
  return <div>{props.user.username}</div>;
}

export default UserStatus;
