import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { convertIntToDateTime } from "./Common";

function FeedbackPage() {
  const [content, setContext] = useState("");
  const [feedbacks, setFeedbacks] = useState([]);

  function getFeedbacks() {
    fetch("/api/v1/get_feedbacks")
      .then((res) => res.json())
      .then((data) => {
        if (data.error) {
          console.log(data.error);
        } else {
          setFeedbacks(data.feedbacks);
        }
      });
  }

  function submitFeedback() {
    fetch("/api/v1/feedback", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        content: content,
      }),
    })
      .then((res) => res.json())
      .then((data) => {
        if (data.error) {
          alert(data.error);
        } else {
          setContext("");
          getFeedbacks();
        }
      });
  }

  useEffect(() => {
    getFeedbacks();
  }, []);

  return (
    <div className="page">
      <h3>Feedback</h3>
      <textarea value={content} onChange={(e) => setContext(e.target.value)} />
      <button onClick={() => submitFeedback()}>Submit</button>
      <div>
        <table>
          <thead>
            <tr>
              <th>User</th>
              <th>Feedback</th>
              <th>Date</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            {feedbacks.map((feedback) => (
              <tr key={feedback._id}>
                <td>
                  <Link to={`/manage/users/${feedback.user.id}`}>
                    @{feedback.user.username}
                  </Link>
                </td>
                <td>{feedback.content}</td>
                <td>{convertIntToDateTime(feedback.time)}</td>
                <td>
                  <button onClick={() => {}}>Delete</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default FeedbackPage;
