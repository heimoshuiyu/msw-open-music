import { useState, useEffect } from "react";
import { useParams } from "react-router";
import ReviewEntry from "./ReviewEntry";

function UserProfile(props) {
  let params = useParams();
  const [reviews, setReviews] = useState([]);

  function getReviews() {
    fetch("/api/v1/get_reviews_by_user", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
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

  useEffect(() => {
    getReviews();
  }, []);

  return (
    <div className="page">
      <h3>User Profile</h3>
      <div>
        <h4>Reviews</h4>
        {reviews.map((review) => (
          <ReviewEntry
            key={review.id}
            review={review}
            user={props.user}
          />
        ))}
      </div>
    </div>
  );
}

export default UserProfile;
