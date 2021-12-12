import { useState } from "react";
import { useParams } from "react-router";

function ReviewPage() {
  let params = useParams();
  const [newReview, setNewReview] = useState("");

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
        }
      });
  }

  return (
    <div className="page">
      <h3>Review Page</h3>
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
