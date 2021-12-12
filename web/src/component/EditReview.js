import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router";

function SingleReview() {
  let params = useParams();
  let navigate = useNavigate();

  const [review, setReview] = useState({
    id: "",
    user_id: "",
    file_id: "",
    content: "",
    created_at: "",
    updated_at: "",
  });

  function refresh() {
    fetch("/api/v1/get_review", {
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
          setReview(data.review);
        }
      });
  }

  function save() {
    fetch("/api/v1/update_review", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        id: parseInt(params.id),
        content: review.content,
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          alert(data.error);
        } else {
          alert("Review updated!");
          navigate(-1);
        }
      });
  }

  useEffect(() => {
    refresh();
  }, []);

  return (
    <div className="page">
      <h3>Edit Review</h3>
      <textarea
        value={review.content}
        onChange={(e) => setReview({ ...review, content: e.target.value })}
      ></textarea>
      <button onClick={() => save()}>Save</button>
    </div>
  );
}

export default SingleReview;
