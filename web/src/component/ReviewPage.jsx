import * as React from 'react';
import { useState, useEffect } from "react";
import { useParams } from "react-router";
import ReviewEntry from "./ReviewEntry";
import { Tr } from "../translate";

function ReviewPage(props) {
  let params = useParams();
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
      <h3>{Tr("Review Page")}</h3>
      <div>
        {reviews.map((review) => (
          <ReviewEntry key={review.id} review={review} user={props.user} />
        ))}
      </div>
      <div>
        <textarea
          value={newReview}
          onChange={(e) => setNewReview(e.target.value)}
        />
        <button onClick={() => submitReview()}>{Tr("Submit")}</button>
      </div>
    </div>
  );
}

export default ReviewPage;
