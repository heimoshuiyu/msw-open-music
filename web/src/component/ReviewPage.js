import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router";
import { Link } from "react-router-dom";
import { convertIntToDateTime } from "./Common";

function ReviewPage(props) {
  let params = useParams();
  let navigate = useNavigate();
  const [newReview, setNewReview] = useState("");
  const [reviews, setReviews] = useState([]);

  function refresh() {
    fetch("/api/v1/get_reviews_on_file", {
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
          setReviews(data.reviews);
        }
      });
  }

  useEffect(() => {
    refresh();
  }, []);

  function submitReview() {
    fetch("/api/v1/insert_review", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        content: newReview,
        file_id: parseInt(params.id),
      }),
    })
      .then((res) => res.json())
      .then((data) => {
        if (data.error) {
          alert(data.error);
        } else {
          setNewReview("");
          refresh();
        }
      });
  }

  return (
    <div className="page">
      <h3>Review Page</h3>
      <div>
        {reviews.map((review) => (
          <div key={review.id}>
            <h4>
              <Link to={`/manage/users/${review.user.id}`}>
                @{review.user.username}
              </Link>{" "}
              wrote on {convertIntToDateTime(review.created_at)}{" "}
            </h4>
            <p>{review.content}</p>
            {(props.user.role === 1 || review.user.id === props.user.id) &&
              props.user.role != 0 && (
                <button
                  onClick={() => {
                    navigate(`/manage/reviews/${review.id}`);
                  }}
                >
                  Edit
                </button>
              )}
          </div>
        ))}
      </div>
      <div>
        <textarea
          value={newReview}
          onChange={(e) => setNewReview(e.target.value)}
        />
        <button onClick={() => submitReview()}>Submit</button>
      </div>
    </div>
  );
}

export default ReviewPage;
